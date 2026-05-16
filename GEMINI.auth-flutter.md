# GEMINI.auth-flutter.md — MANTRA Flutter: Auth API Integration

> Tugas sesi ini: audit kondisi Flutter saat ini, lalu sambungkan auth flow
> ke endpoint backend yang sudah live.
> Jangan ubah UI/widget apapun kecuali diminta eksplisit.

---

## FASE 1 — AUDIT DULU, JANGAN CODING

Sebelum menulis satu baris kode, lakukan hal berikut:

### 1.1 Scan struktur folder
Baca dan pahami struktur `frontend/lib/` secara keseluruhan. Identifikasi:
- Folder apa saja yang ada (`screens/`, `services/`, `providers/`, `models/`, dll)
- State management yang dipakai (Provider, Riverpod, BLoC, GetX, atau belum ada)
- HTTP client yang dipakai (`dio`, `http`, atau belum ada)
- Apakah sudah ada `flutter_secure_storage` di `pubspec.yaml`
- Apakah sudah ada file service/repository untuk auth (`auth_service.dart`, `api_service.dart`, dll)
- Apakah sudah ada model untuk User (`user_model.dart` atau sejenisnya)

### 1.2 Cek `pubspec.yaml`
Catat semua dependency yang relevan dengan networking dan storage.

### 1.3 Temukan screen auth yang sudah ada
Cari file screen untuk:
- Login
- Register
- Splash / onboarding (biasanya cek token di sini)

### 1.4 Buat laporan audit singkat
Setelah scan, **tulis ringkasan** kondisi saat ini sebelum mulai coding:
```
AUDIT RESULT:
- HTTP client: [ada/belum — nama package]
- State management: [nama / belum ada]
- flutter_secure_storage: [ada/belum]
- Auth service: [ada/belum — path file]
- User model: [ada/belum — path file]
- Login screen: [path file]
- Register screen: [path file / tidak ada]
- Splash screen: [path file / tidak ada]
```

---

## FASE 2 — SETUP (jika belum ada)

Lakukan hanya untuk yang **belum ada** berdasarkan hasil audit.

### 2.1 Tambah dependency jika belum ada

Di `pubspec.yaml`, pastikan ada:
```yaml
dependencies:
  dio: ^5.x.x                      # HTTP client
  flutter_secure_storage: ^9.x.x   # simpan token
  # jika pakai state management, sesuaikan
```

Jalankan `flutter pub get` setelah edit.

### 2.2 Buat `ApiClient` jika belum ada

Buat di `lib/core/network/api_client.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiClient {
  static const String baseUrl = String.fromEnvironment(
    'BASE_URL',
    defaultValue: 'http://10.0.2.2:8080/api/v1', // emulator Android
  );

  final Dio _dio;
  final FlutterSecureStorage _storage;

  ApiClient({Dio? dio, FlutterSecureStorage? storage})
      : _dio = dio ?? Dio(BaseOptions(baseUrl: baseUrl)),
        _storage = storage ?? const FlutterSecureStorage() {
    _dio.interceptors.add(_AuthInterceptor(_storage, _dio));
  }

  Dio get dio => _dio;
}

class _AuthInterceptor extends Interceptor {
  final FlutterSecureStorage _storage;
  final Dio _dio;

  _AuthInterceptor(this._storage, this._dio);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.read(key: 'access_token');
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // Coba refresh token
      final refreshed = await _tryRefresh();
      if (refreshed) {
        // Retry request asal
        final token = await _storage.read(key: 'access_token');
        err.requestOptions.headers['Authorization'] = 'Bearer $token';
        final response = await _dio.fetch(err.requestOptions);
        handler.resolve(response);
        return;
      }
      // Refresh gagal — hapus token, redirect ke login
      await _storage.deleteAll();
    }
    handler.next(err);
  }

  Future<bool> _tryRefresh() async {
    try {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) return false;

      final response = await Dio().post(
        '${ApiClient.baseUrl}/auth/refresh',
        data: {'refresh_token': refreshToken},
      );
      final newToken = response.data['data']['access_token'];
      await _storage.write(key: 'access_token', value: newToken);
      return true;
    } catch (_) {
      return false;
    }
  }
}
```

### 2.3 Buat `User` model jika belum ada

Buat di `lib/core/models/user_model.dart`:

```dart
class UserModel {
  final int idUser;
  final String username;
  final String email;
  final String namaLengkap;
  final String role;
  final int? profileId;
  final String? fotoProfil;

  UserModel({
    required this.idUser,
    required this.username,
    required this.email,
    required this.namaLengkap,
    required this.role,
    this.profileId,
    this.fotoProfil,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      idUser: json['id_user'],
      username: json['username'],
      email: json['email'],
      namaLengkap: json['nama_lengkap'],
      role: json['role'],
      profileId: json['profile_id'],
      fotoProfil: json['foto_profil'],
    );
  }
}
```

---

## KONTEKS BACKEND (SUDAH DIKETAHUI — JANGAN DIUBAH)

Hasil audit backend yang sudah selesai dan perlu diketahui Flutter:

- **JWT claims:** `user_id`, `public_id`, `role`, `exp`, `iat`
- **Refresh token:** berupa **random bytes hex** (bukan JWT) — disimpan dan dicocokkan dari DB
- **Response login Flutter:**
  ```json
  {
    "data": {
      "access_token": "...",
      "token_type": "Bearer",
      "expires_in": 900,
      "user": {
        "id_user": 1,
        "username": "...",
        "email": "...",
        "nama_lengkap": "...",
        "role": "customer",
        "profile_id": 10
      }
    }
  }
  ```
- **Response refresh token:** `{ "data": { "access_token": "..." } }`
- Backend deteksi Flutter dari ada-tidaknya header `Authorization: Bearer`

---

## FASE 3 — BUAT AUTH SERVICE

Buat atau update file auth service. Letakkan di path yang konsisten dengan struktur
yang sudah ada (dari hasil audit). Jika belum ada, buat di
`lib/features/auth/services/auth_service.dart`.

```dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../../core/models/user_model.dart';

class AuthService {
  final Dio _dio;
  final FlutterSecureStorage _storage;

  AuthService(this._dio, this._storage);

  // LOGIN
  // Endpoint: POST /api/v1/login
  Future<UserModel> login(String username, String password) async {
    final response = await _dio.post('/login', data: {
      'username': username,
      'password': password,
    });

    final data = response.data['data'];

    // Simpan token ke secure storage
    // access_token = JWT, refresh_token = random bytes hex
    await _storage.write(key: 'access_token', value: data['access_token']);
    await _storage.write(key: 'refresh_token', value: data['refresh_token']);
    await _storage.write(key: 'role', value: data['user']['role']);

    return UserModel.fromJson(data['user']);
  }

  // REGISTER
  // Endpoint: POST /api/v1/register
  Future<void> register({
    required String username,
    required String email,
    required String password,
    required String konfirmasiPassword,
    required String namaLengkap,
    required String noTelp,
  }) async {
    await _dio.post('/register', data: {
      'username': username,
      'email': email,
      'password': password,
      'konfirmasi_password': konfirmasiPassword,
      'nama_lengkap': namaLengkap,
      'no_telp': noTelp,
    });
    // Tidak auto-login setelah register — arahkan ke halaman login
  }

  // LOGOUT
  // Endpoint: POST /api/v1/logout
  Future<void> logout() async {
    final refreshToken = await _storage.read(key: 'refresh_token');
    try {
      await _dio.post('/logout', data: {'refresh_token': refreshToken});
    } catch (_) {
      // Tetap hapus token lokal meskipun request gagal
    } finally {
      await _storage.deleteAll();
    }
  }

  // CEK SESI — dipakai di splash screen
  Future<bool> isLoggedIn() async {
    final token = await _storage.read(key: 'access_token');
    return token != null;
  }

  // AMBIL ROLE DARI STORAGE — untuk routing setelah login
  Future<String?> getSavedRole() async {
    return await _storage.read(key: 'role');
  }
}
```

---

## FASE 4 — SAMBUNGKAN KE SCREEN YANG SUDAH ADA

> Jangan buat screen baru. Sambungkan `AuthService` ke screen yang sudah ada
> dari hasil audit Fase 1.

### 4.1 Login Screen

Temukan fungsi/method yang dipanggil saat tombol login ditekan. Ganti logika
dummy/placeholder dengan:

```dart
// Contoh pola — sesuaikan dengan state management yang dipakai
Future<void> _handleLogin() async {
  setState(() => _isLoading = true);
  try {
    final user = await authService.login(
      _usernameController.text.trim(),
      _passwordController.text,
    );

    // Simpan role untuk routing
    await storage.write(key: 'role', value: user.role);

    // Routing berdasarkan role
    switch (user.role) {
      case 'customer':
        Navigator.pushReplacementNamed(context, '/customer/home');
        break;
      case 'kasir':
        Navigator.pushReplacementNamed(context, '/kasir/home');
        break;
      case 'admin':
        Navigator.pushReplacementNamed(context, '/admin/home');
        break;
    }
  } on DioException catch (e) {
    final code = e.response?.data['error']['code'];
    final message = e.response?.data['message'] ?? 'Login gagal';
    // Tampilkan error ke UI — sesuaikan dengan cara show error yang sudah dipakai
    _showError(message);
  } finally {
    setState(() => _isLoading = false);
  }
}
```

### 4.2 Register Screen (jika ada)

Sambungkan tombol register ke `authService.register(...)`. Setelah berhasil,
arahkan ke login screen dengan pesan sukses.

### 4.3 Splash Screen (jika ada)

Ganti logika cek sesi dengan:

```dart
Future<void> _checkSession() async {
  final isLoggedIn = await authService.isLoggedIn();
  if (!isLoggedIn) {
    Navigator.pushReplacementNamed(context, '/login');
    return;
  }

  final role = await authService.getSavedRole();
  switch (role) {
    case 'customer':
      Navigator.pushReplacementNamed(context, '/customer/home');
      break;
    case 'kasir':
      Navigator.pushReplacementNamed(context, '/kasir/home');
      break;
    case 'admin':
      Navigator.pushReplacementNamed(context, '/admin/home');
      break;
    default:
      Navigator.pushReplacementNamed(context, '/login');
  }
}
```

---

## FASE 5 — ERROR HANDLING STANDAR

Semua error dari API mengikuti format:
```json
{
  "status": "error",
  "message": "Pesan untuk UI",
  "error": { "code": "AUTH_001", "detail": "..." }
}
```

Buat helper jika belum ada di `lib/core/utils/api_error.dart`:

```dart
class ApiError {
  final String code;
  final String message;

  ApiError({required this.code, required this.message});

  factory ApiError.fromDioException(DioException e) {
    final data = e.response?.data;
    return ApiError(
      code: data?['error']?['code'] ?? 'SERVER_001',
      message: data?['message'] ?? 'Terjadi kesalahan, coba lagi',
    );
  }

  // Pesan khusus per kode error auth
  String get userMessage {
    switch (code) {
      case 'AUTH_001': return 'Username atau password salah';
      case 'AUTH_003': return 'Sesi habis, silakan login kembali';
      case 'CONF_001': return 'Username sudah digunakan';
      case 'CONF_002': return 'Email sudah terdaftar';
      case 'VAL_002':  return 'Konfirmasi password tidak cocok';
      case 'VAL_003':  return 'Format email tidak valid';
      default:         return message;
    }
  }
}
```

---

## CHECKLIST SESI INI

- [ ] Audit selesai dan laporan ditulis sebelum coding dimulai
- [ ] `dio` dan `flutter_secure_storage` ada di `pubspec.yaml`
- [ ] `ApiClient` dengan auto-refresh interceptor sudah ada
- [ ] `AuthService` dengan `login`, `register`, `logout`, `isLoggedIn` sudah ada
- [ ] Login screen sudah pakai `AuthService` (bukan dummy)
- [ ] Register screen sudah pakai `AuthService` (jika screen-nya ada)
- [ ] Splash screen sudah cek token dan route berdasarkan role (jika screen-nya ada)
- [ ] Error handling pakai kode error dari backend (`AUTH_001`, `CONF_001`, dll)
- [ ] Token disimpan di `flutter_secure_storage`, bukan `SharedPreferences`
- [ ] Tidak ada token atau credential yang hardcode di kode

---

## YANG TIDAK BOLEH DILAKUKAN

- Jangan buat screen/widget baru — sambungkan ke yang sudah ada
- Jangan simpan token di `SharedPreferences` — harus `flutter_secure_storage`
- Jangan hardcode `BASE_URL` langsung di service — taruh di `ApiClient` saja
- Jangan ubah UI/design/layout apapun
- Jangan asumsikan nama route — cek `MaterialApp` atau router yang sudah ada
- Jangan auto-login setelah register — arahkan ke halaman login
