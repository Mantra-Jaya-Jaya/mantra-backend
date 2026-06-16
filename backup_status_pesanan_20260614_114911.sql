--
-- PostgreSQL database dump
--

\restrict lFqV4I3g6gMjLcWOkldox15LDX5VkhhvDl9dIAgfTh6wQKyn4q7nBIj7kKyNC3h

-- Dumped from database version 18.4
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: alamat; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.alamat (
    id_alamat bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    id_customer bigint,
    nama_penerima text,
    label_alamat text,
    no_telp_penerima text,
    alamat_lengkap text,
    kode_pos text,
    latitude numeric,
    longitude numeric,
    catatan_lokasi text,
    is_utama boolean
);


ALTER TABLE public.alamat OWNER TO postgres;

--
-- Name: alamat_id_alamat_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.alamat_id_alamat_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.alamat_id_alamat_seq OWNER TO postgres;

--
-- Name: alamat_id_alamat_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.alamat_id_alamat_seq OWNED BY public.alamat.id_alamat;


--
-- Name: barang; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.barang (
    id_barang bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    nama_barang text,
    gambar_barang text,
    deskripsi text,
    panjang_barang bigint DEFAULT 0,
    lebar_barang bigint DEFAULT 0,
    tinggi_barang bigint DEFAULT 0,
    id_diskon bigint,
    id_satuan bigint,
    id_kategori bigint
);


ALTER TABLE public.barang OWNER TO postgres;

--
-- Name: barang_id_barang_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.barang_id_barang_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.barang_id_barang_seq OWNER TO postgres;

--
-- Name: barang_id_barang_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.barang_id_barang_seq OWNED BY public.barang.id_barang;


--
-- Name: barcode; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.barcode (
    id_barcode bigint NOT NULL,
    kode_barcode character varying,
    kuantitas bigint,
    id_spesifikasi_barang bigint
);


ALTER TABLE public.barcode OWNER TO postgres;

--
-- Name: barcode_id_barcode_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.barcode_id_barcode_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.barcode_id_barcode_seq OWNER TO postgres;

--
-- Name: barcode_id_barcode_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.barcode_id_barcode_seq OWNED BY public.barcode.id_barcode;


--
-- Name: customer; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customer (
    id_customer bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    no_telp text,
    id_user bigint
);


ALTER TABLE public.customer OWNER TO postgres;

--
-- Name: customer_id_customer_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.customer_id_customer_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.customer_id_customer_seq OWNER TO postgres;

--
-- Name: customer_id_customer_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.customer_id_customer_seq OWNED BY public.customer.id_customer;


--
-- Name: detail_pembayaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detail_pembayaran (
    id_detail_pembayaran bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    id_pembayaran bigint,
    kanal_pembayaran text,
    nomor_va text,
    bill_key text,
    bill_code text,
    nama_bank text,
    merchant_id text,
    qr_code_url text
);


ALTER TABLE public.detail_pembayaran OWNER TO postgres;

--
-- Name: detail_pembayaran_id_detail_pembayaran_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.detail_pembayaran_id_detail_pembayaran_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.detail_pembayaran_id_detail_pembayaran_seq OWNER TO postgres;

--
-- Name: detail_pembayaran_id_detail_pembayaran_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.detail_pembayaran_id_detail_pembayaran_seq OWNED BY public.detail_pembayaran.id_detail_pembayaran;


--
-- Name: detail_pesanan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detail_pesanan (
    id_detail_pesanan bigint NOT NULL,
    jumlah bigint,
    harga_satuan bigint,
    subtotal bigint,
    id_pesanan bigint,
    id_spesifikasi_barang bigint
);


ALTER TABLE public.detail_pesanan OWNER TO postgres;

--
-- Name: detail_pesanan_id_detail_pesanan_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.detail_pesanan_id_detail_pesanan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.detail_pesanan_id_detail_pesanan_seq OWNER TO postgres;

--
-- Name: detail_pesanan_id_detail_pesanan_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.detail_pesanan_id_detail_pesanan_seq OWNED BY public.detail_pesanan.id_detail_pesanan;


--
-- Name: detail_spesifikasi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detail_spesifikasi (
    id_detail_spesifikasi bigint NOT NULL,
    nama_detail_spesifikasi text,
    id_spesifikasi bigint
);


ALTER TABLE public.detail_spesifikasi OWNER TO postgres;

--
-- Name: detail_spesifikasi_id_detail_spesifikasi_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.detail_spesifikasi_id_detail_spesifikasi_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.detail_spesifikasi_id_detail_spesifikasi_seq OWNER TO postgres;

--
-- Name: detail_spesifikasi_id_detail_spesifikasi_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.detail_spesifikasi_id_detail_spesifikasi_seq OWNED BY public.detail_spesifikasi.id_detail_spesifikasi;


--
-- Name: diskon; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.diskon (
    id_diskon bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    nama_diskon text,
    besar_diskon bigint,
    banner_diskon text,
    tgl_mulai date,
    tgl_selesai date
);


ALTER TABLE public.diskon OWNER TO postgres;

--
-- Name: diskon_id_diskon_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.diskon_id_diskon_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.diskon_id_diskon_seq OWNER TO postgres;

--
-- Name: diskon_id_diskon_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.diskon_id_diskon_seq OWNED BY public.diskon.id_diskon;


--
-- Name: ekspedisi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ekspedisi (
    id_ekspedisi bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    nama_ekspedisi text,
    kode_api text,
    logo text,
    deskripsi text,
    is_active boolean DEFAULT true
);


ALTER TABLE public.ekspedisi OWNER TO postgres;

--
-- Name: ekspedisi_id_ekspedisi_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ekspedisi_id_ekspedisi_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ekspedisi_id_ekspedisi_seq OWNER TO postgres;

--
-- Name: ekspedisi_id_ekspedisi_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ekspedisi_id_ekspedisi_seq OWNED BY public.ekspedisi.id_ekspedisi;


--
-- Name: ekspedisi_layanan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ekspedisi_layanan (
    id_ekspedisi_layanan bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    id_ekspedisi bigint,
    nama_layanan text,
    deskripsi text,
    estimasi_min bigint,
    estimasi_max bigint,
    is_active boolean DEFAULT true
);


ALTER TABLE public.ekspedisi_layanan OWNER TO postgres;

--
-- Name: ekspedisi_layanan_id_ekspedisi_layanan_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ekspedisi_layanan_id_ekspedisi_layanan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ekspedisi_layanan_id_ekspedisi_layanan_seq OWNER TO postgres;

--
-- Name: ekspedisi_layanan_id_ekspedisi_layanan_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ekspedisi_layanan_id_ekspedisi_layanan_seq OWNED BY public.ekspedisi_layanan.id_ekspedisi_layanan;


--
-- Name: karyawan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.karyawan (
    id_karyawan bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    no_telp character varying(15) NOT NULL,
    tempat_lahir text,
    tanggal_lahir date,
    jenis_kelamin text,
    alamat text,
    pendidikan_terakhir text,
    nik character varying(16) NOT NULL,
    status text,
    id_user bigint
);


ALTER TABLE public.karyawan OWNER TO postgres;

--
-- Name: karyawan_id_karyawan_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.karyawan_id_karyawan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.karyawan_id_karyawan_seq OWNER TO postgres;

--
-- Name: karyawan_id_karyawan_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.karyawan_id_karyawan_seq OWNED BY public.karyawan.id_karyawan;


--
-- Name: kasir; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kasir (
    id_kasir bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    shift text,
    id_karyawan bigint
);


ALTER TABLE public.kasir OWNER TO postgres;

--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kasir_id_kasir_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kasir_id_kasir_seq OWNER TO postgres;

--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kasir_id_kasir_seq OWNED BY public.kasir.id_kasir;


--
-- Name: kategori; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kategori (
    id_kategori bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    nama_kategori text,
    icon_kategori text
);


ALTER TABLE public.kategori OWNER TO postgres;

--
-- Name: kategori_id_kategori_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kategori_id_kategori_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kategori_id_kategori_seq OWNER TO postgres;

--
-- Name: kategori_id_kategori_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kategori_id_kategori_seq OWNED BY public.kategori.id_kategori;


--
-- Name: keranjang; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.keranjang (
    id_keranjang bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    quantity bigint,
    id_customer bigint,
    id_spesifikasi_barang bigint
);


ALTER TABLE public.keranjang OWNER TO postgres;

--
-- Name: keranjang_id_keranjang_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.keranjang_id_keranjang_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.keranjang_id_keranjang_seq OWNER TO postgres;

--
-- Name: keranjang_id_keranjang_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.keranjang_id_keranjang_seq OWNED BY public.keranjang.id_keranjang;


--
-- Name: kurir; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kurir (
    id_kurir bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    id_karyawan bigint
);


ALTER TABLE public.kurir OWNER TO postgres;

--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.kurir_id_kurir_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.kurir_id_kurir_seq OWNER TO postgres;

--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.kurir_id_kurir_seq OWNED BY public.kurir.id_kurir;


--
-- Name: metode_pembayaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.metode_pembayaran (
    id_metode_pembayaran bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    nama_metode text,
    kode_metode text,
    penyedia text,
    icon text,
    urutan bigint DEFAULT 0,
    is_active boolean DEFAULT true
);


ALTER TABLE public.metode_pembayaran OWNER TO postgres;

--
-- Name: metode_pembayaran_id_metode_pembayaran_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.metode_pembayaran_id_metode_pembayaran_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.metode_pembayaran_id_metode_pembayaran_seq OWNER TO postgres;

--
-- Name: metode_pembayaran_id_metode_pembayaran_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.metode_pembayaran_id_metode_pembayaran_seq OWNED BY public.metode_pembayaran.id_metode_pembayaran;


--
-- Name: notifikasi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifikasi (
    id_notifikasi bigint NOT NULL,
    id_user bigint,
    judul text,
    pesan text,
    status text,
    created_at timestamp with time zone
);


ALTER TABLE public.notifikasi OWNER TO postgres;

--
-- Name: notifikasi_id_notifikasi_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifikasi_id_notifikasi_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notifikasi_id_notifikasi_seq OWNER TO postgres;

--
-- Name: notifikasi_id_notifikasi_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifikasi_id_notifikasi_seq OWNED BY public.notifikasi.id_notifikasi;


--
-- Name: pembayaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pembayaran (
    id_pembayaran bigint NOT NULL,
    order_id_midtrans text,
    payment_type text,
    status_transaksi text,
    fraud_status text,
    id_pesanan bigint,
    id_metode_pembayaran bigint,
    transaksi_midtrans_id text,
    waktu_pembayaran timestamp with time zone,
    total_dibayar bigint DEFAULT 0
);


ALTER TABLE public.pembayaran OWNER TO postgres;

--
-- Name: pembayaran_id_pembayaran_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pembayaran_id_pembayaran_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pembayaran_id_pembayaran_seq OWNER TO postgres;

--
-- Name: pembayaran_id_pembayaran_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pembayaran_id_pembayaran_seq OWNED BY public.pembayaran.id_pembayaran;


--
-- Name: pengantaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pengantaran (
    id_pengantaran bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    waktu_pickup timestamp with time zone,
    waktu_sampai timestamp with time zone,
    last_latitude numeric,
    last_longitude numeric,
    foto_bukti_pengiriman text,
    id_pesanan bigint,
    id_kurir bigint,
    id_status_pengantaran bigint,
    id_ekspedisi bigint
);


ALTER TABLE public.pengantaran OWNER TO postgres;

--
-- Name: pengantaran_id_pengantaran_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pengantaran_id_pengantaran_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pengantaran_id_pengantaran_seq OWNER TO postgres;

--
-- Name: pengantaran_id_pengantaran_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pengantaran_id_pengantaran_seq OWNED BY public.pengantaran.id_pengantaran;


--
-- Name: pesanan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pesanan (
    id_pesanan bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    total_pembayaran bigint,
    tanggal_pesanan timestamp with time zone,
    tipe_pesanan text,
    status_pesanan text,
    id_customer bigint,
    id_kasir bigint,
    id_alamat bigint,
    id_ekspedisi bigint,
    id_layanan_ekspedisi bigint,
    ongkos_kirim bigint DEFAULT 0,
    catatan text
);


ALTER TABLE public.pesanan OWNER TO postgres;

--
-- Name: pesanan_id_pesanan_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pesanan_id_pesanan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pesanan_id_pesanan_seq OWNER TO postgres;

--
-- Name: pesanan_id_pesanan_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pesanan_id_pesanan_seq OWNED BY public.pesanan.id_pesanan;


--
-- Name: refresh_token; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.refresh_token (
    id_token bigint NOT NULL,
    token text,
    expires_at timestamp with time zone,
    created_at timestamp with time zone,
    revoked_at timestamp with time zone,
    id_user bigint
);


ALTER TABLE public.refresh_token OWNER TO postgres;

--
-- Name: refresh_token_id_token_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.refresh_token_id_token_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_token_id_token_seq OWNER TO postgres;

--
-- Name: refresh_token_id_token_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.refresh_token_id_token_seq OWNED BY public.refresh_token.id_token;


--
-- Name: role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role (
    id_role bigint NOT NULL,
    nama_role text
);


ALTER TABLE public.role OWNER TO postgres;

--
-- Name: role_id_role_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.role_id_role_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.role_id_role_seq OWNER TO postgres;

--
-- Name: role_id_role_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.role_id_role_seq OWNED BY public.role.id_role;


--
-- Name: satuan; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.satuan (
    id_satuan bigint NOT NULL,
    nama_satuan text
);


ALTER TABLE public.satuan OWNER TO postgres;

--
-- Name: satuan_id_satuan_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.satuan_id_satuan_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.satuan_id_satuan_seq OWNER TO postgres;

--
-- Name: satuan_id_satuan_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.satuan_id_satuan_seq OWNED BY public.satuan.id_satuan;


--
-- Name: spesifikasi; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.spesifikasi (
    id_spesifikasi bigint NOT NULL,
    nama_spesifikasi text
);


ALTER TABLE public.spesifikasi OWNER TO postgres;

--
-- Name: spesifikasi_barang; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.spesifikasi_barang (
    id_spesifikasi_barang bigint NOT NULL,
    jumlah bigint,
    harga_barang bigint,
    berat_barang bigint DEFAULT 0,
    id_barang bigint,
    id_detail_spesifikasi bigint
);


ALTER TABLE public.spesifikasi_barang OWNER TO postgres;

--
-- Name: spesifikasi_barang_id_spesifikasi_barang_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.spesifikasi_barang_id_spesifikasi_barang_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.spesifikasi_barang_id_spesifikasi_barang_seq OWNER TO postgres;

--
-- Name: spesifikasi_barang_id_spesifikasi_barang_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.spesifikasi_barang_id_spesifikasi_barang_seq OWNED BY public.spesifikasi_barang.id_spesifikasi_barang;


--
-- Name: spesifikasi_id_spesifikasi_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.spesifikasi_id_spesifikasi_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.spesifikasi_id_spesifikasi_seq OWNER TO postgres;

--
-- Name: spesifikasi_id_spesifikasi_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.spesifikasi_id_spesifikasi_seq OWNED BY public.spesifikasi.id_spesifikasi;


--
-- Name: status_pengantaran; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.status_pengantaran (
    id_status_pengantaran bigint NOT NULL,
    nama_status text
);


ALTER TABLE public.status_pengantaran OWNER TO postgres;

--
-- Name: status_pengantaran_id_status_pengantaran_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.status_pengantaran_id_status_pengantaran_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.status_pengantaran_id_status_pengantaran_seq OWNER TO postgres;

--
-- Name: status_pengantaran_id_status_pengantaran_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.status_pengantaran_id_status_pengantaran_seq OWNED BY public.status_pengantaran.id_status_pengantaran;


--
-- Name: stok_opname; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stok_opname (
    id_stok_opname bigint NOT NULL,
    harga_beli bigint,
    status boolean,
    jumlah_stok bigint,
    keterangan text,
    tanggal timestamp with time zone,
    id_spesifikasi_barang bigint
);


ALTER TABLE public.stok_opname OWNER TO postgres;

--
-- Name: stok_opname_id_stok_opname_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stok_opname_id_stok_opname_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stok_opname_id_stok_opname_seq OWNER TO postgres;

--
-- Name: stok_opname_id_stok_opname_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stok_opname_id_stok_opname_seq OWNED BY public.stok_opname.id_stok_opname;


--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id_user bigint NOT NULL,
    public_id uuid DEFAULT gen_random_uuid(),
    username text,
    email text,
    password text,
    nama_lengkap text,
    foto_profil text,
    id_role bigint
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- Name: user_id_user_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id_user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_id_user_seq OWNER TO postgres;

--
-- Name: user_id_user_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_id_user_seq OWNED BY public."user".id_user;


--
-- Name: alamat id_alamat; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.alamat ALTER COLUMN id_alamat SET DEFAULT nextval('public.alamat_id_alamat_seq'::regclass);


--
-- Name: barang id_barang; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barang ALTER COLUMN id_barang SET DEFAULT nextval('public.barang_id_barang_seq'::regclass);


--
-- Name: barcode id_barcode; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barcode ALTER COLUMN id_barcode SET DEFAULT nextval('public.barcode_id_barcode_seq'::regclass);


--
-- Name: customer id_customer; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer ALTER COLUMN id_customer SET DEFAULT nextval('public.customer_id_customer_seq'::regclass);


--
-- Name: detail_pembayaran id_detail_pembayaran; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pembayaran ALTER COLUMN id_detail_pembayaran SET DEFAULT nextval('public.detail_pembayaran_id_detail_pembayaran_seq'::regclass);


--
-- Name: detail_pesanan id_detail_pesanan; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pesanan ALTER COLUMN id_detail_pesanan SET DEFAULT nextval('public.detail_pesanan_id_detail_pesanan_seq'::regclass);


--
-- Name: detail_spesifikasi id_detail_spesifikasi; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_spesifikasi ALTER COLUMN id_detail_spesifikasi SET DEFAULT nextval('public.detail_spesifikasi_id_detail_spesifikasi_seq'::regclass);


--
-- Name: diskon id_diskon; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.diskon ALTER COLUMN id_diskon SET DEFAULT nextval('public.diskon_id_diskon_seq'::regclass);


--
-- Name: ekspedisi id_ekspedisi; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ekspedisi ALTER COLUMN id_ekspedisi SET DEFAULT nextval('public.ekspedisi_id_ekspedisi_seq'::regclass);


--
-- Name: ekspedisi_layanan id_ekspedisi_layanan; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ekspedisi_layanan ALTER COLUMN id_ekspedisi_layanan SET DEFAULT nextval('public.ekspedisi_layanan_id_ekspedisi_layanan_seq'::regclass);


--
-- Name: karyawan id_karyawan; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan ALTER COLUMN id_karyawan SET DEFAULT nextval('public.karyawan_id_karyawan_seq'::regclass);


--
-- Name: kasir id_kasir; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir ALTER COLUMN id_kasir SET DEFAULT nextval('public.kasir_id_kasir_seq'::regclass);


--
-- Name: kategori id_kategori; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kategori ALTER COLUMN id_kategori SET DEFAULT nextval('public.kategori_id_kategori_seq'::regclass);


--
-- Name: keranjang id_keranjang; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keranjang ALTER COLUMN id_keranjang SET DEFAULT nextval('public.keranjang_id_keranjang_seq'::regclass);


--
-- Name: kurir id_kurir; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir ALTER COLUMN id_kurir SET DEFAULT nextval('public.kurir_id_kurir_seq'::regclass);


--
-- Name: metode_pembayaran id_metode_pembayaran; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.metode_pembayaran ALTER COLUMN id_metode_pembayaran SET DEFAULT nextval('public.metode_pembayaran_id_metode_pembayaran_seq'::regclass);


--
-- Name: notifikasi id_notifikasi; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi ALTER COLUMN id_notifikasi SET DEFAULT nextval('public.notifikasi_id_notifikasi_seq'::regclass);


--
-- Name: pembayaran id_pembayaran; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pembayaran ALTER COLUMN id_pembayaran SET DEFAULT nextval('public.pembayaran_id_pembayaran_seq'::regclass);


--
-- Name: pengantaran id_pengantaran; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran ALTER COLUMN id_pengantaran SET DEFAULT nextval('public.pengantaran_id_pengantaran_seq'::regclass);


--
-- Name: pesanan id_pesanan; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan ALTER COLUMN id_pesanan SET DEFAULT nextval('public.pesanan_id_pesanan_seq'::regclass);


--
-- Name: refresh_token id_token; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_token ALTER COLUMN id_token SET DEFAULT nextval('public.refresh_token_id_token_seq'::regclass);


--
-- Name: role id_role; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role ALTER COLUMN id_role SET DEFAULT nextval('public.role_id_role_seq'::regclass);


--
-- Name: satuan id_satuan; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.satuan ALTER COLUMN id_satuan SET DEFAULT nextval('public.satuan_id_satuan_seq'::regclass);


--
-- Name: spesifikasi id_spesifikasi; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi ALTER COLUMN id_spesifikasi SET DEFAULT nextval('public.spesifikasi_id_spesifikasi_seq'::regclass);


--
-- Name: spesifikasi_barang id_spesifikasi_barang; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi_barang ALTER COLUMN id_spesifikasi_barang SET DEFAULT nextval('public.spesifikasi_barang_id_spesifikasi_barang_seq'::regclass);


--
-- Name: status_pengantaran id_status_pengantaran; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_pengantaran ALTER COLUMN id_status_pengantaran SET DEFAULT nextval('public.status_pengantaran_id_status_pengantaran_seq'::regclass);


--
-- Name: stok_opname id_stok_opname; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stok_opname ALTER COLUMN id_stok_opname SET DEFAULT nextval('public.stok_opname_id_stok_opname_seq'::regclass);


--
-- Name: user id_user; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user" ALTER COLUMN id_user SET DEFAULT nextval('public.user_id_user_seq'::regclass);


--
-- Data for Name: alamat; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.alamat (id_alamat, public_id, id_customer, nama_penerima, label_alamat, no_telp_penerima, alamat_lengkap, kode_pos, latitude, longitude, catatan_lokasi, is_utama) FROM stdin;
4	7a6fa854-0763-499e-8a84-925c24b31fc9	3	Rajaba Hamim	Kost	084370056578	Jl. Banjarsari Selatan, Tembalang, Kota Semarang		-7.05141	110.438125	Pagar hitam, samping warung burjo	t
5	1266b8b6-5801-4e2f-bbf6-478a614f4666	3	Rajaba Hamim	Rumah	087754156937	Kecamatan Selogiri, Kabupaten Wonogiri		-7.816667	110.916667	Rumah cat hijau dekat pertigaan balai desa	f
1	8318e207-7dd3-4f43-b2db-491acb6cbb59	1	Rajaba Hamim	Kost	088558366847	Jl. Banjarsari Selatan, Tembalang, Kota Semarang		-7.05141	110.438125	Pagar hitam, samping warung burjo	f
2	6be71a78-6dfd-4019-9f0b-4bd37f72e7e2	1	Rajaba Hamim	Rumah	082489625497	Kecamatan Selogiri, Kabupaten Wonogiri		-7.816667	110.916667	Rumah cat hijau dekat pertigaan balai desa	f
3	c2adb034-175b-42ec-ab37-c66a01f679c5	1	John Doe	Rumah	08123456789	Jl. Contoh No. 1, Jakarta		-6.2	106.8	Depan masjid	f
6	890f64f8-1766-44f6-8aa6-1964eb308747	1	John Doe	Rumah	08123456789	Jl. Contoh No. 1, Jakarta		-6.2	106.8	Depan masjid	t
\.


--
-- Data for Name: barang; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.barang (id_barang, public_id, nama_barang, gambar_barang, deskripsi, panjang_barang, lebar_barang, tinggi_barang, id_diskon, id_satuan, id_kategori) FROM stdin;
2	863b83d1-b9fd-44da-8e5a-b1409eb34f75	Smartphone Samsung A55	https://picsum.photos/seed/196/400/400	Invite review for the problem in Cleveland.	0	0	0	1	1	1
3	c4e045e8-3db6-40b6-957c-5a31f55835c4	TWS Earphone Bluetooth	https://picsum.photos/seed/527/400/400	Track company over time weekly.	0	0	0	1	1	1
4	10569328-9976-4c84-9ed4-1f5f97db931f	Sepatu Lari Nike Air	https://picsum.photos/seed/45/400/400	Create a fallback for child.	0	0	0	1	1	2
5	9b555322-bd59-4c41-af61-bb2e5b0c0e9f	Kaos Polos Premium	https://picsum.photos/seed/627/400/400	Phew! Ship the week now!	0	0	0	1	1	2
6	7207b0c8-b9cb-4bb6-9943-818d2e6bf7b8	Jaket Hoodie Fleece	https://picsum.photos/seed/121/400/400	Ruthlessly remove dead group.	0	0	0	1	1	2
7	00047a48-765c-4287-9fd6-94cf2c7a7eed	Mie Instan Box 40pcs	https://picsum.photos/seed/231/400/400	Eventually, the way wake violently.	0	0	0	1	1	3
8	e2581ff8-4b51-4e8a-afb7-a0fd02aad12a	Kopi Sachet Box	https://picsum.photos/seed/129/400/400	Guard eye with sensible limits.	0	0	0	1	1	3
9	67898de0-2c04-4fe7-8c21-31629ed680bc	Susu UHT Full Cream 1L	https://picsum.photos/seed/705/400/400	Explicitly name the work before you question it.	0	0	0	1	1	3
10	bf457891-00d1-429a-a820-c9770e85a3ce	Vitamin C 1000mg 30 tablet	https://picsum.photos/seed/644/400/400	Sample life at 5s intervals.	0	0	0	1	1	4
11	01d38324-e27a-4781-9ff3-7b00a788f691	Masker KN95 Box isi 20	https://picsum.photos/seed/912/400/400	Share the decision record for the year.	0	0	0	1	1	4
12	aa58f630-11cf-4f7e-be71-a691a1122244	Hand Sanitizer 500ml	https://picsum.photos/seed/421/400/400	Track place over time nightly.	0	0	0	1	1	4
13	68295608-4e9d-4140-8bcb-742674140057	Dumbbell Set 5kg	https://picsum.photos/seed/591/400/400	Surface risks around the child too.	0	0	0	1	1	5
14	69d3b393-b05f-4dcf-9e08-8682f8bae493	Matras Yoga Anti Slip	https://picsum.photos/seed/226/400/400	Guard life with sensible limits.	0	0	0	1	1	5
15	9b497273-2508-401a-8226-10f949ee72ce	Raket Badminton Carbon	https://picsum.photos/seed/592/400/400	Retire outdated thing each quarter.	0	0	0	1	1	5
1	4df3846b-6592-4821-9773-ffb3c3bf3b8a	Laptop Gaming ASUS ROG	https://picsum.photos/seed/332/400/400	Invite review for the fact in Fremont.	0	0	0	6	1	1
\.


--
-- Data for Name: barcode; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.barcode (id_barcode, kode_barcode, kuantitas, id_spesifikasi_barang) FROM stdin;
5	279267908647	1	3
6	140500815790	12	3
7	417999074846	1	4
8	254918873029	12	4
9	905642664543	1	5
10	559222255995	12	5
11	816764325868	1	6
12	700737940712	12	6
13	877891939850	1	7
14	767505174952	12	7
15	293508414683	1	8
16	684691826387	12	8
17	512175378278	1	9
18	142238121741	12	9
19	135923657250	1	10
20	904388601946	12	10
21	902088460996	1	11
22	672872306279	12	11
23	844527883601	1	12
24	273270064477	12	12
25	206577983988	1	13
26	843216564724	12	13
27	253189146865	1	14
28	491577961612	12	14
29	405948495542	1	15
30	386797290436	12	15
31	759162522093	1	16
32	577556965312	12	16
33	841961602797	1	17
34	336008026036	12	17
35	857450978993	1	18
36	833470106189	12	18
37	318324998102	1	19
38	852232462701	12	19
39	552191989740	1	20
40	402537299681	12	20
41	651594586368	1	21
42	744605785256	12	21
43	531882165891	1	22
44	875995515999	12	22
45	358916458286	1	23
46	908066163356	12	23
47	416656688474	1	24
48	665268882339	12	24
49	934277895981	1	25
50	294267850332	12	25
51	501897867863	1	26
52	195895939958	12	26
53	204804158057	1	27
54	262561649057	12	27
55	654520819709	1	28
56	251671225407	12	28
57	834474731493	1	29
58	763575038742	12	29
59	784667599378	1	30
60	667604214553	12	30
67	227383890177	1	3
68	747076393730	12	3
69	400846059775	1	4
70	664418805913	12	4
71	455318429888	1	5
72	601694333346	12	5
73	684959678505	1	6
74	933479471294	12	6
75	466177847778	1	7
76	766827831681	12	7
77	173740927352	1	8
78	965244365059	12	8
79	356634321492	1	9
80	859467507646	12	9
81	119477146091	1	10
82	654256428301	12	10
83	736198565383	1	11
84	111879444307	12	11
85	630630338851	1	12
86	888117294988	12	12
87	444719952121	1	13
88	342049409678	12	13
89	720906755334	1	14
90	507087338215	12	14
91	111448889415	1	15
92	488254897574	12	15
93	417331052251	1	16
94	684964219064	12	16
95	225070467596	1	17
96	550806320411	12	17
97	353688415932	1	18
98	700241625084	12	18
99	654894561097	1	19
100	466482288204	12	19
101	867063715709	1	20
102	870609355910	12	20
103	365948036274	1	21
104	543725005754	12	21
105	372815525301	1	22
106	220058791123	12	22
107	866280055726	1	23
108	771083078869	12	23
109	396940680719	1	24
110	575097186761	12	24
111	497703719816	1	25
112	367558808215	12	25
113	208272642464	1	26
114	778871041460	12	26
115	892926500195	1	27
116	124263607770	12	27
117	562338689022	1	28
118	319202004956	12	28
119	207570550837	1	29
120	523683891021	12	29
121	157415875348	1	30
122	325665520775	12	30
127	536417562508	1	3
128	399179100742	12	3
129	411915931367	1	4
130	165336362897	12	4
131	789505569911	1	5
132	659391024927	12	5
133	951137567122	1	6
134	246392471134	12	6
135	824138773490	1	7
136	377762414012	12	7
137	794791597149	1	8
138	193938184978	12	8
139	285481257148	1	9
140	510268466079	12	9
141	282198015780	1	10
142	376413975645	12	10
143	971721479490	1	11
144	476365025268	12	11
145	559254086148	1	12
146	578200498006	12	12
147	120049286748	1	13
148	844168322074	12	13
149	311675648970	1	14
150	251294293448	12	14
151	579700660963	1	15
152	852866129439	12	15
153	188600794897	1	16
154	435906458136	12	16
155	713028452429	1	17
156	798999969101	12	17
157	771031385270	1	18
158	457853038691	12	18
159	692240813860	1	19
160	693777637069	12	19
161	498462507488	1	20
162	567072109400	12	20
163	712256011941	1	21
164	968467649722	12	21
165	979823071676	1	22
166	306644720231	12	22
167	269612227100	1	23
168	511879840480	12	23
169	926359257860	1	24
170	711937633206	12	24
171	580694351743	1	25
172	439968160449	12	25
173	458162584787	1	26
174	966881115686	12	26
175	444384470445	1	27
176	326644140835	12	27
177	436261473001	1	28
178	931869132048	12	28
179	104015758409	1	29
180	390901251919	12	29
181	381137864659	1	30
182	680561488159	12	30
189	898632570116	1	3
190	485361381036	12	3
191	555842866808	1	4
192	493571089693	12	4
193	432781505786	1	5
194	338147280417	12	5
195	800914738762	1	6
196	318689805361	12	6
197	950779863215	1	7
198	943224487583	12	7
199	198260168214	1	8
200	647464405365	12	8
201	420385003805	1	9
202	808580901069	12	9
203	518687610917	1	10
204	665224906411	12	10
205	925598865955	1	11
206	681456713679	12	11
207	375706348150	1	12
208	948184963031	12	12
209	248580680996	1	13
210	693215727923	12	13
211	455530417784	1	14
212	893820630754	12	14
213	576036015955	1	15
214	131301211293	12	15
215	815972000667	1	16
216	896023198595	12	16
217	146524448668	1	17
218	872676438533	12	17
219	571644383254	1	18
220	162990617512	12	18
221	639571611497	1	19
222	447011068126	12	19
223	738361815169	1	20
224	921321578480	12	20
225	769006438721	1	21
226	434609882067	12	21
227	132797800522	1	22
228	265622657419	12	22
229	357710435482	1	23
230	282351532255	12	23
231	540075220538	1	24
232	387734148675	12	24
233	579707450977	1	25
234	181569808980	12	25
235	268187034946	1	26
236	221797822193	12	26
237	337315470905	1	27
238	531877239847	12	27
239	532669214000	1	28
240	326311866145	12	28
241	872460116181	1	29
242	310298173432	12	29
243	662402266094	1	30
244	108650019642	12	30
249	804389911792	1	3
250	625645692305	12	3
251	653285129232	1	4
252	972294682131	12	4
253	389845950795	1	5
254	118694685878	12	5
255	314829359388	1	6
256	596606256851	12	6
257	443357323685	1	7
258	905104034868	12	7
259	347094862593	1	8
260	868459864161	12	8
261	101605534956	1	9
262	893825191699	12	9
263	484038853210	1	10
264	935445075963	12	10
265	182708371519	1	11
266	629502295310	12	11
267	804797199699	1	12
268	342839037066	12	12
269	622126736663	1	13
270	961398510837	12	13
271	465630131828	1	14
272	323279356766	12	14
273	356991437485	1	15
274	415905445287	12	15
275	662959962717	1	16
276	349164593845	12	16
277	990794700197	1	17
278	879690776726	12	17
279	711856287191	1	18
280	193209104914	12	18
281	462103670152	1	19
282	544668797052	12	19
283	716645875012	1	20
284	411738322373	12	20
285	108253967271	1	21
286	275902209856	12	21
287	899789107429	1	22
288	342929371935	12	22
289	308394565579	1	23
290	418373060251	12	23
291	853916256992	1	24
292	406297889196	12	24
293	230541361080	1	25
294	146098788893	12	25
295	162700580355	1	26
296	395548222670	12	26
297	286331623394	1	27
298	696728020931	12	27
299	939399088316	1	28
300	319934473835	12	28
301	114383117110	1	29
302	558776214616	12	29
303	216492546870	1	30
304	841699547545	12	30
305	853547272053	1	1
306	131905771719	12	1
307	578554109591	1	1
308	712752870688	12	1
309	197295681272	1	1
310	455222749457	12	1
311	743120630619	1	1
312	207569897451	12	1
313	454251457956	1	1
314	536845041486	12	1
315	590419244306	1	2
316	911686167589	12	2
317	971452779340	1	2
318	332133653172	12	2
319	526273035479	1	2
320	878314040470	12	2
321	143397736643	1	2
322	837475141691	12	2
323	582549233629	1	2
324	356195571360	12	2
\.


--
-- Data for Name: customer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.customer (id_customer, public_id, no_telp, id_user) FROM stdin;
2	5ef40515-4f6f-4bb8-98de-eac432c44434	081233	6
3	fe20634f-18d3-4d37-ae5c-4abe8f852661	085080283720	7
4	fe9d802c-9cf9-44bd-a208-707741ca9301	0812997	9
1	a3752eea-c812-466c-b1d0-8f084b4f0a6e	08123456780	4
5	e90cda9f-48c5-4508-9927-c5d4c23550b0	0812601	10
\.


--
-- Data for Name: detail_pembayaran; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.detail_pembayaran (id_detail_pembayaran, public_id, id_pembayaran, kanal_pembayaran, nomor_va, bill_key, bill_code, nama_bank, merchant_id, qr_code_url) FROM stdin;
\.


--
-- Data for Name: detail_pesanan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.detail_pesanan (id_detail_pesanan, jumlah, harga_satuan, subtotal, id_pesanan, id_spesifikasi_barang) FROM stdin;
1	5	193009	965045	1	18
2	4	423793	1695172	1	27
3	5	240763	1203815	1	29
4	2	626578	1253156	2	7
5	4	263949	1055796	2	8
6	1	14816845	14816845	3	2
7	2	423793	847586	3	27
8	2	30543	61086	3	21
9	1	626578	626578	4	7
10	5	775768	3878840	4	12
11	4	775768	3103072	5	12
12	3	40359	121077	5	13
13	1	54842	54842	6	24
14	3	775529	2326587	6	9
15	4	112185	448740	6	28
16	3	10002883	30008649	7	5
17	4	38911	155644	7	23
18	3	21363	64089	7	22
19	1	240763	240763	7	29
20	1	105736	105736	8	14
21	3	162448	487344	8	16
22	2	3229255	6458510	8	6
23	4	10002883	40011532	9	5
24	5	67111	335555	9	20
25	1	14097752	14097752	9	3
26	1	97462	97462	10	19
27	3	12861980	38585940	10	1
28	2	67111	134222	11	20
29	5	172895	864475	11	30
30	5	105736	528680	12	14
31	2	112185	224370	12	28
32	2	162448	324896	12	16
33	3	775768	2327304	12	12
34	5	40359	201795	13	13
35	1	38911	38911	13	23
36	5	455963	2279815	13	25
37	1	52841	52841	14	15
38	4	10002883	40011532	14	5
39	2	52841	105682	15	15
40	5	775768	3878840	15	12
41	5	30543	152715	15	21
42	3	10002883	30008649	15	5
43	1	263949	263949	16	8
44	5	14816845	74084225	16	2
45	1	14097752	14097752	16	3
46	4	775768	3103072	17	12
47	1	21363	21363	17	22
48	2	54842	109684	18	24
49	3	775768	2327304	18	12
50	2	172895	345790	18	30
51	2	12861980	25723960	19	1
52	1	263949	263949	19	8
53	3	162448	487344	19	16
54	1	775768	775768	20	12
55	2	14097752	28195504	20	3
56	1	30543	30543	21	21
57	1	455250	455250	21	10
58	1	455963	455963	22	25
59	1	775768	775768	22	12
60	3	455250	1365750	22	10
61	1	97462	97462	23	19
62	3	196418	589254	23	17
63	4	241842	967368	23	26
64	4	162448	649792	24	16
65	3	14816845	44450535	24	2
66	1	241842	241842	25	26
67	2	10002883	20005766	25	5
68	2	97462	194924	26	19
69	1	105736	105736	26	14
70	4	240763	963052	26	29
71	1	12861980	12861980	27	1
72	1	775768	775768	27	12
73	1	196418	196418	27	17
74	1	14097752	14097752	28	3
75	5	14816845	74084225	28	2
76	1	240763	240763	29	29
77	4	38911	155644	29	23
78	3	423793	1271379	29	27
79	1	455963	455963	29	25
80	4	30543	122172	30	21
81	4	52841	211364	30	15
82	2	626578	1253156	31	7
83	5	112185	560925	31	28
84	4	775768	3103072	31	12
85	5	38911	194555	32	23
86	2	776868	1553736	32	11
87	1	455250	455250	32	10
88	3	162448	487344	32	16
89	2	776868	1553736	33	11
90	3	67111	201333	33	20
91	2	14816845	29633690	34	2
92	3	263949	791847	34	8
93	3	193009	579027	35	18
94	3	240763	722289	35	29
95	4	241842	967368	35	26
96	3	14816845	44450535	36	2
97	5	52841	264205	36	15
98	5	776868	3884340	37	11
99	4	455963	1823852	37	25
100	2	30543	61086	38	21
101	3	5588932	16766796	38	4
102	5	626578	3132890	39	7
103	3	54842	164526	39	24
104	3	423793	1271379	40	27
105	1	775768	775768	40	12
106	1	775768	775768	41	12
107	3	67111	201333	41	20
108	5	21363	106815	42	22
109	2	241842	483684	42	26
110	5	30543	152715	42	21
111	5	105736	528680	42	14
112	1	14097752	14097752	43	3
113	1	38911	38911	43	23
114	2	776868	1553736	44	11
115	4	455963	1823852	44	25
116	3	10002883	30008649	45	5
117	5	423793	2118965	45	27
118	3	105736	317208	45	14
119	5	38911	194555	46	23
120	5	30543	152715	46	21
121	5	67111	335555	46	20
122	5	105736	528680	46	14
123	1	3229255	3229255	47	6
124	5	14097752	70488760	47	3
125	2	263949	527898	47	8
126	3	776868	2330604	47	11
127	3	30543	91629	48	21
128	4	196418	785672	48	17
129	2	455250	910500	48	10
130	3	5588932	16766796	49	4
131	2	455963	911926	49	25
132	4	775529	3102116	50	9
133	2	3229255	6458510	50	6
134	3	10002883	30008649	50	5
135	1	14816845	14816845	51	2
136	3	40359	121077	51	13
137	2	196418	392836	51	17
138	2	626578	1253156	51	7
139	2	97462	194924	52	19
140	4	14816845	59267380	52	2
141	2	455963	911926	52	25
142	3	162448	487344	52	16
143	5	263949	1319745	53	8
144	5	112185	560925	53	28
145	2	626578	1253156	53	7
146	3	423793	1271379	54	27
147	5	3229255	16146275	55	6
148	1	626578	626578	55	7
149	3	67111	201333	55	20
150	2	10002883	20005766	56	5
151	4	240763	963052	56	29
152	4	162448	649792	57	16
153	3	193009	579027	57	18
154	1	21363	21363	58	22
155	5	776868	3884340	58	11
156	3	10002883	30008649	59	5
157	4	105736	422944	59	14
158	3	162448	487344	59	16
159	3	172895	518685	60	30
160	3	263949	791847	60	8
161	4	196418	785672	60	17
162	3	775768	2327304	61	12
163	4	40359	161436	61	13
164	2	14816845	29633690	61	2
165	3	196418	589254	62	17
166	4	241842	967368	62	26
167	1	40359	40359	63	13
168	1	67111	67111	63	20
169	4	423793	1695172	64	27
170	2	38911	77822	64	23
171	2	97462	194924	64	19
172	3	162448	487344	65	16
173	1	455250	455250	65	10
174	1	67111	67111	65	20
175	1	775529	775529	66	9
176	2	54842	109684	66	24
177	2	3229255	6458510	66	6
178	4	626578	2506312	66	7
179	2	193009	386018	67	18
180	1	97462	97462	67	19
181	1	30543	30543	67	21
182	2	5588932	11177864	68	4
183	1	97462	97462	68	19
184	1	21363	21363	68	22
185	5	172895	864475	68	30
186	3	626578	1879734	69	7
187	5	162448	812240	69	16
188	5	196418	982090	70	17
189	4	52841	211364	70	15
190	3	112185	336555	70	28
191	2	40359	80718	71	13
192	2	14097752	28195504	71	3
193	3	40359	121077	72	13
194	3	5588932	16766796	72	4
195	1	193009	193009	72	18
196	5	196418	982090	73	17
197	5	21363	106815	73	22
198	2	38911	77822	73	23
199	1	40359	40359	73	13
200	4	12861980	51447920	74	1
201	3	67111	201333	74	20
202	4	54842	219368	74	24
203	2	14816845	29633690	74	2
204	5	12861980	64309900	75	1
205	4	14816845	59267380	75	2
206	3	14097752	42293256	76	3
207	5	105736	528680	76	14
208	2	3229255	6458510	76	6
209	2	776868	1553736	76	11
210	4	67111	268444	77	20
211	1	30543	30543	77	21
212	5	776868	3884340	77	11
213	5	10002883	50014415	77	5
214	3	196418	589254	78	17
215	3	12861980	38585940	78	1
216	1	5588932	5588932	79	4
217	1	172895	172895	79	30
218	1	263949	263949	79	8
219	4	162448	649792	79	16
220	1	54842	54842	80	24
221	4	12861980	51447920	80	1
222	3	240763	722289	81	29
223	2	5588932	11177864	81	4
224	1	21363	21363	82	22
225	1	172895	172895	82	30
226	4	775768	3103072	82	12
227	1	14097752	14097752	83	3
228	3	38911	116733	83	23
229	2	67111	134222	83	20
230	1	776868	776868	84	11
231	4	40359	161436	84	13
232	3	5588932	16766796	84	4
233	2	21363	42726	84	22
234	5	241842	1209210	85	26
235	4	776868	3107472	85	11
236	4	423793	1695172	85	27
237	4	775529	3102116	86	9
238	2	455250	910500	86	10
239	2	3229255	6458510	87	6
240	2	67111	134222	87	20
241	4	423793	1695172	88	27
242	1	14816845	14816845	88	2
243	5	97462	487310	89	19
244	3	52841	158523	89	15
245	2	240763	481526	89	29
246	3	30543	91629	90	21
247	4	455250	1821000	90	10
248	1	196418	196418	90	17
249	4	97462	389848	91	19
250	3	12861980	38585940	91	1
251	2	193009	386018	92	18
252	5	54842	274210	92	24
253	1	52841	52841	92	15
254	4	3229255	12917020	92	6
255	3	172895	518685	93	30
256	4	10002883	40011532	93	5
257	1	455963	455963	94	25
258	3	263949	791847	94	8
259	3	105736	317208	95	14
260	3	423793	1271379	95	27
261	5	193009	965045	95	18
262	5	240763	1203815	95	29
263	3	263949	791847	96	8
264	4	38911	155644	96	23
265	5	196418	982090	96	17
266	5	40359	201795	97	13
267	5	193009	965045	97	18
268	2	3229255	6458510	98	6
269	3	5588932	16766796	98	4
270	4	455963	1823852	99	25
271	5	14097752	70488760	99	3
272	3	193009	579027	99	18
273	2	5588932	11177864	99	4
274	1	105736	105736	100	14
275	5	38911	194555	100	23
276	5	196418	982090	101	17
277	1	12861980	12861980	101	1
278	1	97462	97462	102	19
279	1	455250	455250	102	10
280	1	240763	240763	102	29
281	2	38911	77822	103	23
282	2	52841	105682	103	15
283	3	775768	2327304	103	12
284	3	775768	2327304	104	12
285	1	14816845	14816845	104	2
286	5	423793	2118965	105	27
287	5	5588932	27944660	105	4
288	2	97462	194924	105	19
289	4	5588932	22355728	106	4
290	3	54842	164526	106	24
291	5	21363	106815	107	22
292	1	54842	54842	107	24
293	3	52841	158523	107	15
294	1	5588932	5588932	107	4
295	2	193009	386018	108	18
296	1	12861980	12861980	108	1
297	4	455963	1823852	108	25
298	5	54842	274210	109	24
299	5	196418	982090	109	17
300	2	5588932	11177864	110	4
301	4	14097752	56391008	110	3
302	5	97462	487310	111	19
303	2	112185	224370	111	28
304	2	14816845	29633690	112	2
305	1	172895	172895	112	30
306	1	112185	112185	113	28
307	4	162448	649792	113	16
308	5	67111	335555	113	20
309	5	776868	3884340	113	11
310	4	21363	85452	114	22
311	4	40359	161436	114	13
312	5	776868	3884340	115	11
313	5	196418	982090	115	17
314	3	30543	91629	115	21
315	4	97462	389848	116	19
316	2	193009	386018	116	18
317	4	455250	1821000	116	10
318	2	423793	847586	117	27
319	3	775529	2326587	117	9
320	1	10002883	10002883	117	5
321	1	14097752	14097752	117	3
322	2	776868	1553736	118	11
323	1	196418	196418	118	17
324	3	162448	487344	118	16
325	5	67111	335555	119	20
326	3	97462	292386	119	19
327	1	3229255	3229255	119	6
328	2	775768	1551536	119	12
329	3	97462	292386	120	19
330	3	21363	64089	120	22
331	5	172895	864475	121	30
332	4	455963	1823852	121	25
333	5	67111	335555	122	20
334	3	21363	64089	122	22
335	2	776868	1553736	122	11
336	2	97462	194924	122	19
337	3	112185	336555	123	28
338	5	38911	194555	123	23
339	3	14816845	44450535	123	2
340	3	38911	116733	124	23
341	5	626578	3132890	124	7
342	4	5588932	22355728	124	4
343	1	38911	38911	125	23
344	4	196418	785672	125	17
345	2	105736	211472	125	14
346	5	112185	560925	125	28
347	2	196418	392836	126	17
348	2	423793	847586	127	27
349	2	10002883	20005766	127	5
350	2	172895	345790	127	30
351	1	193009	193009	128	18
352	1	196418	196418	128	17
353	5	3229255	16146275	129	6
354	2	52841	105682	129	15
355	5	776868	3884340	129	11
356	1	775768	775768	129	12
357	2	14816845	29633690	130	2
358	1	172895	172895	130	30
359	1	423793	423793	130	27
360	3	775529	2326587	131	9
361	2	423793	847586	131	27
362	5	30543	152715	132	21
363	4	775529	3102116	132	9
364	1	162448	162448	132	16
365	3	263949	791847	132	8
366	1	455963	455963	133	25
367	2	775768	1551536	133	12
368	5	40359	201795	134	13
369	1	10002883	10002883	134	5
370	1	97462	97462	135	19
371	3	626578	1879734	135	7
372	3	263949	791847	136	8
373	2	196418	392836	136	17
374	5	12861980	64309900	136	1
375	4	455963	1823852	137	25
376	4	3229255	12917020	137	6
377	3	112185	336555	137	28
378	3	775529	2326587	137	9
379	4	455250	1821000	138	10
380	4	52841	211364	138	15
381	2	241842	483684	138	26
382	4	196418	785672	138	17
383	4	172895	691580	139	30
384	3	241842	725526	139	26
385	4	263949	1055796	139	8
386	2	3229255	6458510	140	6
387	4	5588932	22355728	140	4
388	1	455963	455963	140	25
389	3	40359	121077	141	13
390	4	776868	3107472	141	11
391	4	193009	772036	142	18
392	5	626578	3132890	142	7
393	1	196418	196418	143	17
394	2	455963	911926	143	25
395	2	67111	134222	143	20
396	3	455250	1365750	144	10
397	5	12861980	64309900	144	1
398	3	30543	91629	145	21
399	3	263949	791847	145	8
400	1	112185	112185	145	28
401	1	162448	162448	145	16
402	3	775768	2327304	146	12
403	3	776868	2330604	146	11
404	3	196418	589254	146	17
405	4	455963	1823852	146	25
406	5	172895	864475	147	30
407	4	30543	122172	147	21
408	4	5588932	22355728	147	4
409	2	162448	324896	147	16
410	1	626578	626578	148	7
411	1	776868	776868	148	11
412	2	193009	386018	149	18
413	5	423793	2118965	149	27
414	5	241842	1209210	150	26
415	4	196418	785672	150	17
416	3	423793	1271379	151	27
417	3	162448	487344	151	16
418	2	241842	483684	152	26
419	1	626578	626578	152	7
420	5	196418	982090	152	17
421	4	263949	1055796	153	8
422	2	38911	77822	153	23
423	5	423793	2118965	153	27
424	4	54842	219368	153	24
425	2	3229255	6458510	154	6
426	2	21363	42726	154	22
427	5	12861980	64309900	155	1
428	3	162448	487344	155	16
429	1	240763	240763	155	29
430	3	775529	2326587	155	9
431	3	21363	64089	156	22
432	5	14097752	70488760	156	3
433	2	52841	105682	156	15
434	5	10002883	50014415	157	5
435	5	38911	194555	157	23
436	2	423793	847586	157	27
437	3	775768	2327304	157	12
438	1	455250	455250	158	10
439	3	52841	158523	158	15
440	2	162448	324896	159	16
441	2	455250	910500	159	10
442	4	105736	422944	159	14
443	5	3229255	16146275	159	6
444	2	240763	481526	160	29
445	5	423793	2118965	160	27
446	4	196418	785672	160	17
447	1	38911	38911	161	23
448	2	241842	483684	161	26
449	1	455250	455250	161	10
450	3	5588932	16766796	161	4
451	2	5588932	11177864	162	4
452	5	12861980	64309900	162	1
453	2	240763	481526	162	29
454	5	626578	3132890	162	7
455	2	263949	527898	163	8
456	2	105736	211472	163	14
457	1	12861980	12861980	164	1
458	2	240763	481526	164	29
459	5	172895	864475	164	30
460	2	112185	224370	165	28
461	4	105736	422944	165	14
462	2	263949	527898	165	8
463	5	54842	274210	165	24
464	3	775768	2327304	166	12
465	4	775529	3102116	166	9
466	1	240763	240763	166	29
467	5	193009	965045	167	18
468	5	38911	194555	167	23
469	5	3229255	16146275	168	6
470	1	97462	97462	168	19
471	3	626578	1879734	168	7
472	4	455250	1821000	169	10
473	4	40359	161436	169	13
474	3	240763	722289	169	29
475	4	12861980	51447920	170	1
476	5	54842	274210	170	24
477	3	455963	1367889	171	25
478	5	775768	3878840	171	12
479	4	423793	1695172	172	27
480	1	10002883	10002883	172	5
481	4	67111	268444	172	20
482	3	112185	336555	172	28
483	3	10002883	30008649	173	5
484	1	193009	193009	173	18
485	4	30543	122172	173	21
486	3	423793	1271379	174	27
487	5	105736	528680	174	14
488	5	10002883	50014415	174	5
489	2	30543	61086	174	21
490	2	54842	109684	175	24
491	3	105736	317208	175	14
492	5	196418	982090	175	17
493	4	112185	448740	175	28
494	1	626578	626578	176	7
495	3	455963	1367889	176	25
496	3	40359	121077	177	13
497	5	775529	3877645	177	9
498	1	196418	196418	177	17
499	5	21363	106815	177	22
500	4	626578	2506312	178	7
501	1	14097752	14097752	178	3
502	1	30543	30543	179	21
503	4	3229255	12917020	179	6
504	5	423793	2118965	179	27
505	3	626578	1879734	179	7
506	1	21363	21363	180	22
507	5	455250	2276250	180	10
508	4	54842	219368	180	24
509	5	112185	560925	181	28
510	5	240763	1203815	181	29
511	1	67111	67111	181	20
512	3	423793	1271379	182	27
513	2	240763	481526	182	29
514	4	193009	772036	183	18
515	3	263949	791847	183	8
516	4	67111	268444	183	20
517	4	14816845	59267380	183	2
518	3	455963	1367889	184	25
519	1	162448	162448	184	16
520	2	12861980	25723960	184	1
521	2	263949	527898	185	8
522	4	776868	3107472	185	11
523	4	172895	691580	186	30
524	5	54842	274210	186	24
525	3	172895	518685	187	30
526	3	162448	487344	187	16
527	1	626578	626578	187	7
528	1	112185	112185	187	28
529	1	97462	97462	188	19
530	3	775529	2326587	188	9
531	2	172895	345790	189	30
532	4	14097752	56391008	189	3
533	5	240763	1203815	189	29
534	5	423793	2118965	189	27
535	5	193009	965045	190	18
536	5	105736	528680	190	14
537	2	776868	1553736	191	11
538	4	54842	219368	191	24
539	3	775768	2327304	192	12
540	1	240763	240763	192	29
541	3	40359	121077	192	13
542	1	3229255	3229255	192	6
543	5	67111	335555	193	20
544	4	14816845	59267380	193	2
545	4	775529	3102116	193	9
546	1	775529	775529	194	9
547	4	97462	389848	194	19
548	1	196418	196418	194	17
549	1	775768	775768	194	12
550	2	776868	1553736	195	11
551	5	21363	106815	195	22
552	3	455963	1367889	196	25
553	3	3229255	9687765	196	6
554	5	14097752	70488760	197	3
555	4	40359	161436	197	13
556	4	38911	155644	198	23
557	2	241842	483684	198	26
558	1	196418	196418	199	17
559	4	14097752	56391008	199	3
560	3	38911	116733	200	23
561	3	21363	64089	200	22
562	2	97462	194924	200	19
563	3	52841	158523	200	15
564	5	162448	812240	201	16
565	2	30543	61086	201	21
566	2	52841	105682	201	15
567	3	776868	2330604	202	11
568	5	52841	264205	202	15
569	5	30543	152715	203	21
570	2	21363	42726	203	22
571	3	14816845	44450535	203	2
572	5	626578	3132890	203	7
573	3	112185	336555	204	28
574	5	54842	274210	204	24
575	1	172895	172895	204	30
576	5	21363	106815	204	22
577	3	105736	317208	205	14
578	2	775768	1551536	205	12
579	1	455250	455250	206	10
580	5	240763	1203815	206	29
581	2	30543	61086	206	21
582	2	241842	483684	206	26
583	3	455963	1367889	207	25
584	5	14816845	74084225	207	2
585	2	105736	211472	208	14
586	4	52841	211364	208	15
587	3	241842	725526	208	26
588	2	240763	481526	209	29
589	3	263949	791847	209	8
590	1	775529	775529	209	9
591	3	52841	158523	209	15
592	2	40359	80718	210	13
593	4	455963	1823852	210	25
594	3	423793	1271379	211	27
595	4	67111	268444	211	20
596	4	112185	448740	211	28
597	4	54842	219368	211	24
598	1	10002883	10002883	212	5
599	4	423793	1695172	212	27
600	4	12861980	51447920	213	1
601	4	455250	1821000	213	10
602	2	112185	224370	213	28
603	3	196418	589254	214	17
604	5	105736	528680	214	14
605	2	30543	61086	215	21
606	4	423793	1695172	215	27
607	4	455250	1821000	216	10
608	5	14816845	74084225	216	2
609	5	263949	1319745	217	8
610	3	14816845	44450535	217	2
611	1	30543	30543	217	21
612	4	241842	967368	217	26
613	3	21363	64089	218	22
614	5	3229255	16146275	218	6
615	2	193009	386018	218	18
616	4	626578	2506312	219	7
617	2	97462	194924	219	19
618	1	30543	30543	219	21
619	3	162448	487344	219	16
620	4	52841	211364	220	15
621	1	112185	112185	220	28
622	2	38911	77822	221	23
623	5	14097752	70488760	221	3
624	5	263949	1319745	222	8
625	2	40359	80718	222	13
626	1	172895	172895	222	30
627	5	52841	264205	222	15
628	5	162448	812240	223	16
629	1	105736	105736	223	14
630	4	172895	691580	224	30
631	1	193009	193009	224	18
632	4	193009	772036	225	18
633	2	97462	194924	225	19
634	1	14097752	14097752	225	3
635	1	626578	626578	225	7
636	2	12594319	25188638	226	2
637	2	11983090	23966180	226	3
638	4	10932683	43730732	226	1
639	2	10932683	21865366	227	1
640	2	10932683	21865366	228	1
641	2	12594319	25188638	229	2
642	2	11983090	23966180	229	3
643	4	10932683	43730732	229	1
644	2	11983090	23966180	230	3
645	3	11318543	33955629	230	1
646	4	13038824	52155296	230	2
\.


--
-- Data for Name: detail_spesifikasi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.detail_spesifikasi (id_detail_spesifikasi, nama_detail_spesifikasi, id_spesifikasi) FROM stdin;
1	Merah	1
2	Kuning	1
3	Hijau	1
4	36	2
5	37	2
6	38	2
7	39	2
8	40	2
9	41	2
10	42	2
11	43	2
12	44	2
13	45	2
14	Large	2
\.


--
-- Data for Name: diskon; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.diskon (id_diskon, public_id, nama_diskon, besar_diskon, banner_diskon, tgl_mulai, tgl_selesai) FROM stdin;
1	1602f95a-a04b-488d-9723-4089144ac373	Flash Sale Elektronik	15	https://placehold.co/800x300?text=Flash+Sale+Elektronik	2026-06-09	2026-07-09
3	9006c39d-f786-42ec-81dc-6c0f4d5a2714	Diskon Member Baru	10	https://placehold.co/800x300?text=Diskon+Member+Baru	2026-06-09	2026-08-08
6	df8aeaca-0465-4f57-9a9c-71885d5f1b08	coba1	12	https://storage.mantra.web.id/mantra-storage/diskon/2026/06/4334d64b-13b7-4720-896f-2f2cadb99572.png	2026-06-10	2026-06-13
\.


--
-- Data for Name: ekspedisi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ekspedisi (id_ekspedisi, public_id, nama_ekspedisi, kode_api, logo, deskripsi, is_active) FROM stdin;
1	96f0f8fd-1542-4e8f-a3a8-cc35d07a67cf	SPEX Express	spex			t
2	c1766cfe-7513-4549-8968-aea15ab544b0	JNE	jne			t
3	b1c00ea3-eddc-458e-b015-4c774e2e7121	J&T Express	jnt			t
4	18e96492-6ed0-4d5b-9ae0-eda895186ce4	SiCepat	sicepat			t
5	26647bea-5c3b-418c-b5e9-6cced4440b5e	Anteraja	anteraja			t
6	08d61bc3-65f6-4572-a8b5-2fe266c2c0b2	Ninja Xpress	ninja			t
\.


--
-- Data for Name: ekspedisi_layanan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ekspedisi_layanan (id_ekspedisi_layanan, public_id, id_ekspedisi, nama_layanan, deskripsi, estimasi_min, estimasi_max, is_active) FROM stdin;
1	c72d282d-9602-4d9b-a841-f06b01cd6fca	1	REG	Reguler	2	4	t
2	70ff4156-4c48-434b-b28e-8b54b000c620	2	REG	Reguler	2	3	t
3	025288e3-ef0e-4b5a-bbbe-485060b23de0	2	YES	Yakin Esok Sampai	1	1	t
4	0d665160-6133-40c0-b81d-ecc08510b909	2	OKE	Ongkos Kirim Ekonomis	3	5	t
5	ecaba9ec-05a2-4330-a358-22057e81a4e8	3	EZ	Economy	2	4	t
6	a315a0d2-0fae-4eab-9e56-e233613e3e42	3	REG	Reguler	2	3	t
7	89f76e83-116b-4050-a21f-3ba1f3962378	4	REG	Reguler	1	2	t
8	76784898-7526-43c5-ab50-ed6e1709ce98	4	BEST	Besok Sampai Tujuan	1	1	t
9	b6ff7d1f-fa2b-4a55-9835-b2fb3917d1e3	5	REG	Reguler	2	4	t
10	aa42af60-ae55-4e6b-a901-a82cbab16da7	5	ND	Next Day	1	1	t
11	177ed21e-90bf-4485-8b11-54704e72565e	6	REG	Reguler	2	4	t
\.


--
-- Data for Name: karyawan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.karyawan (id_karyawan, public_id, no_telp, tempat_lahir, tanggal_lahir, jenis_kelamin, alamat, pendidikan_terakhir, nik, status, id_user) FROM stdin;
4	de680e34-9e91-4d95-8f85-f1737e803d33	0812750	Jakarta	1995-01-01	L	Jl. Contoh No. 1	SMA	417	Aktif	8
2	71fc50d1-e64d-4718-a384-bf2d983347f0	082882147802	Wonogiri	2003-10-10	Perempuan	Kecamatan Selogiri, Kabupaten Wonogiri	SMA Negeri 1 Wonogiri	3312428115140565	Aktif	3
1	a8ec160f-fa81-4a44-a212-3211f6dbfe3b	087244293531	Semarang	2004-05-15	Perempuan	Jl. Prof. Sudarto, Tembalang, Kota Semarang	D3 Teknik Komputer	3374959884945274	Aktif	2
3	870bf334-6b38-4bf1-a218-b14484363e73	0812521	Jakarta	1995-01-01	L	Jl. Contoh No. 1	SMA	627	Aktif	5
\.


--
-- Data for Name: kasir; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kasir (id_kasir, public_id, shift, id_karyawan) FROM stdin;
3	650e434a-7fe2-40ed-ae66-dff7da888fff	Pagi	4
1	248789ce-4ea7-4c8d-8a7c-e181bd768ea2	Pagi	1
2	8aed10ac-86e0-47c6-aa9d-d3865629c032	Pagi	3
\.


--
-- Data for Name: kategori; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kategori (id_kategori, public_id, nama_kategori, icon_kategori) FROM stdin;
1	538f32b6-c550-4118-8ea0-2cf407285c19	Elektronik	https://placehold.co/64x64?text=Elektronik
2	d6f902f2-d284-42df-b577-a0d5702f4440	Fashion	https://placehold.co/64x64?text=Fashion
3	6cac4118-0f32-4f37-8e47-067823e94ed7	Makanan & Minuman	https://placehold.co/64x64?text=Makanan+%26+Minuman
4	52302882-7a05-4c15-888c-04b7e36908c6	Kesehatan	https://placehold.co/64x64?text=Kesehatan
5	0513e6f9-9bb4-44b6-8f6a-7c0fe57f8925	Olahraga	https://placehold.co/64x64?text=Olahraga
6	3b6490f5-36d7-4c45-98c1-fde00e9caaad	Peralatan Rumah	https://placehold.co/64x64?text=Peralatan+Rumah
7	3c5d32f4-9d7c-4d26-85ba-10a704b15e87	Buku & Alat Tulis	https://placehold.co/64x64?text=Buku+%26+Alat+Tulis
8	2499bcd4-9a87-49ab-979d-f8c052952436	Kecantikan	https://placehold.co/64x64?text=Kecantikan
11	a6e0d884-871b-468f-8f43-cbecf77695ea	Minuman	
\.


--
-- Data for Name: keranjang; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.keranjang (id_keranjang, public_id, quantity, id_customer, id_spesifikasi_barang) FROM stdin;
\.


--
-- Data for Name: kurir; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.kurir (id_kurir, public_id, id_karyawan) FROM stdin;
1	605052d6-5a4e-4a37-b77e-7b46f658ed28	2
\.


--
-- Data for Name: metode_pembayaran; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.metode_pembayaran (id_metode_pembayaran, public_id, nama_metode, kode_metode, penyedia, icon, urutan, is_active) FROM stdin;
1	78ed793b-8ff4-4009-8363-1ef8ca4ccd74	Cash	cash	internal		1	t
2	8f94d740-b3e9-41d1-825b-c610ca5fd12e	QRIS	qris	midtrans		2	t
3	7c63515c-0530-4306-8c42-2f1da7fa368f	Virtual Account	va	midtrans		3	t
4	efd85050-ec14-4abf-8ee3-6a077c003cdf	E-Wallet	ewallet	midtrans		4	t
5	95e51170-f511-4858-83b3-884a113068b5	COD (Bayar di Tempat)	cod	internal		5	t
\.


--
-- Data for Name: notifikasi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifikasi (id_notifikasi, id_user, judul, pesan, status, created_at) FROM stdin;
1	1	Laporan Harian	Laporan penjualan hari ini sudah tersedia.	unread	2026-06-09 20:23:33.006902+07
2	1	Peringatan Sistem	Perlu pengecekan stok opname bulan ini.	unread	2026-06-09 20:23:33.012727+07
3	1	Karyawan Baru	Ada pendaftaran kasir baru yang menunggu persetujuan.	unread	2026-06-09 20:23:33.019199+07
4	2	Stok Menipis	Beberapa barang hampir habis, segera lakukan restock.	unread	2026-06-09 20:23:33.024842+07
5	2	Transaksi Baru	Ada pesanan online baru yang perlu dikonfirmasi.	unread	2026-06-09 20:23:33.030029+07
6	3	Tugas Pickup	Ada pesanan baru yang harus diambil di toko.	unread	2026-06-09 20:23:33.035278+07
7	3	Rute Terupdate	Perhatikan rute pengantaran karena ada penutupan jalan.	unread	2026-06-09 20:23:33.040733+07
8	4	Pesanan Diterima	Pesanan Anda sedang diproses oleh kasir.	unread	2026-06-09 20:23:33.046186+07
9	4	Pesanan Dikirim	Kurir sedang menuju ke alamat Anda.	unread	2026-06-09 20:23:33.051111+07
10	4	Promo Spesial	Dapatkan diskon 50% untuk produk baru!	unread	2026-06-09 20:23:33.056274+07
\.


--
-- Data for Name: pembayaran; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pembayaran (id_pembayaran, order_id_midtrans, payment_type, status_transaksi, fraud_status, id_pesanan, id_metode_pembayaran, transaksi_midtrans_id, waktu_pembayaran, total_dibayar) FROM stdin;
1	MANTRA-2-1781011410613874269	bank_transfer	settlement	accept	2	\N		\N	0
2	MANTRA-3-1781011410621239226	bank_transfer	settlement	accept	3	\N		\N	0
3	MANTRA-4-1781011410626692681	qris	settlement	accept	4	\N		\N	0
4		cash	settlement	accept	5	\N		\N	0
5	MANTRA-6-1781011410637621263	bank_transfer	settlement	accept	6	\N		\N	0
6		cash	settlement	accept	7	\N		\N	0
7		cash	settlement	accept	8	\N		\N	0
8	MANTRA-9-1781011410655928674	gopay	settlement	accept	9	\N		\N	0
9		cash	settlement	accept	10	\N		\N	0
10	MANTRA-11-1781011410667715242	gopay	settlement	accept	11	\N		\N	0
11		cash	settlement	accept	12	\N		\N	0
12	MANTRA-13-1781011410678393331	gopay	pending	accept	13	\N		\N	0
13		cash	settlement	accept	14	\N		\N	0
14		cash	settlement	accept	15	\N		\N	0
15		cash	cancel	accept	16	\N		\N	0
16	MANTRA-17-1781011410702184599	gopay	settlement	accept	17	\N		\N	0
17		cash	settlement	accept	18	\N		\N	0
18		cash	settlement	accept	19	\N		\N	0
19		cash	settlement	accept	20	\N		\N	0
20	MANTRA-21-1781011410726123062	bank_transfer	settlement	accept	21	\N		\N	0
21		cash	settlement	accept	22	\N		\N	0
22		cash	settlement	accept	24	\N		\N	0
23		cash	settlement	accept	25	\N		\N	0
24	MANTRA-26-1781011410746769796	qris	settlement	accept	26	\N		\N	0
25		cash	settlement	accept	27	\N		\N	0
26	MANTRA-28-1781011410757657563	bank_transfer	settlement	accept	28	\N		\N	0
27		cash	settlement	accept	29	\N		\N	0
28		cash	settlement	accept	30	\N		\N	0
29	MANTRA-31-1781011410774785370	gopay	settlement	accept	31	\N		\N	0
30		cash	settlement	accept	32	\N		\N	0
31	MANTRA-33-1781011410787151059	qris	settlement	accept	33	\N		\N	0
32	MANTRA-34-1781011410792745426	qris	settlement	accept	34	\N		\N	0
33	MANTRA-35-1781011410798158141	bank_transfer	settlement	accept	35	\N		\N	0
34	MANTRA-36-1781011410803317065	bank_transfer	settlement	accept	36	\N		\N	0
35	MANTRA-37-1781011410809061824	bank_transfer	settlement	accept	37	\N		\N	0
36	MANTRA-38-1781011410815276824	gopay	settlement	accept	38	\N		\N	0
37		cash	settlement	accept	39	\N		\N	0
38		cash	settlement	accept	40	\N		\N	0
39		cash	settlement	accept	41	\N		\N	0
40	MANTRA-42-1781011410845021682	gopay	settlement	accept	42	\N		\N	0
41		cash	settlement	accept	43	\N		\N	0
42		cash	settlement	accept	44	\N		\N	0
43	MANTRA-45-1781011410867728149	gopay	settlement	accept	45	\N		\N	0
44	MANTRA-46-1781011410875373258	gopay	settlement	accept	46	\N		\N	0
45	MANTRA-47-1781011410882504061	bank_transfer	settlement	accept	47	\N		\N	0
46		cash	settlement	accept	48	\N		\N	0
47	MANTRA-49-1781011410898839204	qris	cancel	accept	49	\N		\N	0
48		cash	settlement	accept	50	\N		\N	0
49		cash	settlement	accept	51	\N		\N	0
50		cash	settlement	accept	52	\N		\N	0
51		cash	settlement	accept	53	\N		\N	0
52	MANTRA-54-1781011410933815269	gopay	settlement	accept	54	\N		\N	0
53	MANTRA-55-1781011410939909827	gopay	settlement	accept	55	\N		\N	0
54	MANTRA-57-1781011410947073934	gopay	settlement	accept	57	\N		\N	0
55	MANTRA-58-1781011410953329535	qris	settlement	accept	58	\N		\N	0
56	MANTRA-59-1781011410960660050	qris	settlement	accept	59	\N		\N	0
57		cash	settlement	accept	60	\N		\N	0
58		cash	settlement	accept	61	\N		\N	0
59		cash	settlement	accept	62	\N		\N	0
60		cash	settlement	accept	63	\N		\N	0
61	MANTRA-65-1781011410995812296	qris	settlement	accept	65	\N		\N	0
62	MANTRA-66-1781011411002475541	qris	settlement	accept	66	\N		\N	0
63		cash	settlement	accept	67	\N		\N	0
64	MANTRA-68-1781011411015377481	qris	pending	accept	68	\N		\N	0
65	MANTRA-69-1781011411021399311	gopay	settlement	accept	69	\N		\N	0
66	MANTRA-70-1781011411027937000	gopay	settlement	accept	70	\N		\N	0
67	MANTRA-71-1781011411033729475	gopay	settlement	accept	71	\N		\N	0
68	MANTRA-72-1781011411041427040	qris	settlement	accept	72	\N		\N	0
69		cash	settlement	accept	73	\N		\N	0
70	MANTRA-74-1781011411055988461	bank_transfer	settlement	accept	74	\N		\N	0
71	MANTRA-75-1781011411063277075	qris	settlement	accept	75	\N		\N	0
72	MANTRA-76-1781011411071257438	qris	settlement	accept	76	\N		\N	0
73		cash	settlement	accept	77	\N		\N	0
74	MANTRA-78-1781011411083493296	bank_transfer	settlement	accept	78	\N		\N	0
75	MANTRA-79-1781011411090233567	gopay	settlement	accept	79	\N		\N	0
76		cash	settlement	accept	80	\N		\N	0
77	MANTRA-81-1781011411105112445	qris	settlement	accept	81	\N		\N	0
78	MANTRA-82-1781011411113097307	qris	settlement	accept	82	\N		\N	0
79	MANTRA-83-1781011411119311101	qris	settlement	accept	83	\N		\N	0
80		cash	settlement	accept	84	\N		\N	0
81		cash	settlement	accept	85	\N		\N	0
82	MANTRA-86-1781011411139391644	qris	settlement	accept	86	\N		\N	0
83	MANTRA-87-1781011411145538011	gopay	pending	accept	87	\N		\N	0
84	MANTRA-88-1781011411151686876	qris	pending	accept	88	\N		\N	0
85	MANTRA-89-1781011411158062913	bank_transfer	settlement	accept	89	\N		\N	0
86	MANTRA-90-1781011411165124076	gopay	pending	accept	90	\N		\N	0
87	MANTRA-91-1781011411171329559	bank_transfer	settlement	accept	91	\N		\N	0
88		cash	settlement	accept	92	\N		\N	0
89		cash	settlement	accept	93	\N		\N	0
90		cash	settlement	accept	94	\N		\N	0
91		cash	settlement	accept	95	\N		\N	0
92		cash	settlement	accept	96	\N		\N	0
93	MANTRA-98-1781011411210946280	gopay	settlement	accept	98	\N		\N	0
94	MANTRA-99-1781011411217059093	qris	settlement	accept	99	\N		\N	0
95		cash	settlement	accept	100	\N		\N	0
96	MANTRA-101-1781011411229959575	qris	settlement	accept	101	\N		\N	0
97	MANTRA-102-1781011411236121410	qris	settlement	accept	102	\N		\N	0
98	MANTRA-103-1781011411242682452	qris	settlement	accept	103	\N		\N	0
99		cash	settlement	accept	104	\N		\N	0
100	MANTRA-105-1781011411255018518	bank_transfer	settlement	accept	105	\N		\N	0
101		cash	cancel	accept	106	\N		\N	0
102	MANTRA-107-1781011411268903300	qris	pending	accept	107	\N		\N	0
103	MANTRA-108-1781011411275186945	qris	settlement	accept	108	\N		\N	0
104		cash	settlement	accept	109	\N		\N	0
105		cash	settlement	accept	110	\N		\N	0
106	MANTRA-111-1781011411295561853	bank_transfer	settlement	accept	111	\N		\N	0
107	MANTRA-112-1781011411303359056	bank_transfer	settlement	accept	112	\N		\N	0
108	MANTRA-113-1781011411309739568	qris	settlement	accept	113	\N		\N	0
109		cash	settlement	accept	114	\N		\N	0
110	MANTRA-115-1781011411322324928	gopay	settlement	accept	115	\N		\N	0
111	MANTRA-116-1781011411328598922	qris	settlement	accept	116	\N		\N	0
112	MANTRA-117-1781011411334188806	bank_transfer	settlement	accept	117	\N		\N	0
113	MANTRA-118-1781011411341019159	gopay	settlement	accept	118	\N		\N	0
114	MANTRA-119-1781011411348249796	bank_transfer	settlement	accept	119	\N		\N	0
115	MANTRA-120-1781011411354012530	gopay	settlement	accept	120	\N		\N	0
116		cash	settlement	accept	121	\N		\N	0
117	MANTRA-122-1781011411368380467	qris	pending	accept	122	\N		\N	0
118	MANTRA-123-1781011411374975650	gopay	settlement	accept	123	\N		\N	0
119		cash	settlement	accept	124	\N		\N	0
120	MANTRA-125-1781011411388028322	qris	settlement	accept	125	\N		\N	0
121	MANTRA-127-1781011411393880857	bank_transfer	settlement	accept	127	\N		\N	0
122	MANTRA-128-1781011411401444806	bank_transfer	settlement	accept	128	\N		\N	0
123		cash	settlement	accept	129	\N		\N	0
124	MANTRA-130-1781011411414568358	qris	settlement	accept	130	\N		\N	0
125	MANTRA-131-1781011411421225575	qris	settlement	accept	131	\N		\N	0
126	MANTRA-132-1781011411427692897	gopay	pending	accept	132	\N		\N	0
127	MANTRA-133-1781011411433655041	qris	settlement	accept	133	\N		\N	0
128	MANTRA-134-1781011411440525484	qris	settlement	accept	134	\N		\N	0
129		cash	settlement	accept	135	\N		\N	0
130		cash	settlement	accept	136	\N		\N	0
131		cash	cancel	accept	137	\N		\N	0
132		cash	settlement	accept	138	\N		\N	0
133	MANTRA-139-1781011411475546439	bank_transfer	pending	accept	139	\N		\N	0
134	MANTRA-140-1781011411482336530	qris	pending	accept	140	\N		\N	0
135		cash	settlement	accept	141	\N		\N	0
136	MANTRA-142-1781011411495673468	gopay	settlement	accept	142	\N		\N	0
137		cash	settlement	accept	143	\N		\N	0
138	MANTRA-144-1781011411508627476	qris	settlement	accept	144	\N		\N	0
139	MANTRA-145-1781011411514715848	bank_transfer	settlement	accept	145	\N		\N	0
140		cash	settlement	accept	146	\N		\N	0
141	MANTRA-147-1781011411528061056	bank_transfer	settlement	accept	147	\N		\N	0
142	MANTRA-148-1781011411535075655	qris	settlement	accept	148	\N		\N	0
143	MANTRA-149-1781011411541266181	qris	settlement	accept	149	\N		\N	0
144	MANTRA-150-1781011411549009847	bank_transfer	pending	accept	150	\N		\N	0
145	MANTRA-151-1781011411555105960	bank_transfer	settlement	accept	151	\N		\N	0
146		cash	settlement	accept	152	\N		\N	0
147		cash	settlement	accept	153	\N		\N	0
148	MANTRA-154-1781011411577001635	bank_transfer	settlement	accept	154	\N		\N	0
149		cash	settlement	accept	155	\N		\N	0
150	MANTRA-156-1781011411589722733	bank_transfer	settlement	accept	156	\N		\N	0
151		cash	settlement	accept	157	\N		\N	0
152	MANTRA-159-1781011411603056044	qris	settlement	accept	159	\N		\N	0
153	MANTRA-160-1781011411611282163	qris	settlement	accept	160	\N		\N	0
154	MANTRA-161-1781011411617933609	bank_transfer	settlement	accept	161	\N		\N	0
155		cash	cancel	accept	162	\N		\N	0
156		cash	settlement	accept	163	\N		\N	0
157	MANTRA-164-1781011411638113525	bank_transfer	settlement	accept	164	\N		\N	0
158	MANTRA-165-1781011411644034551	bank_transfer	pending	accept	165	\N		\N	0
159		cash	settlement	accept	166	\N		\N	0
160	MANTRA-167-1781011411658609254	bank_transfer	settlement	accept	167	\N		\N	0
161		cash	settlement	accept	168	\N		\N	0
162		cash	settlement	accept	169	\N		\N	0
163	MANTRA-170-1781011411678983976	bank_transfer	settlement	accept	170	\N		\N	0
164	MANTRA-171-1781011411685139688	bank_transfer	settlement	accept	171	\N		\N	0
165	MANTRA-172-1781011411691970927	bank_transfer	settlement	accept	172	\N		\N	0
166		cash	settlement	accept	173	\N		\N	0
167	MANTRA-174-1781011411704337414	bank_transfer	settlement	accept	174	\N		\N	0
168		cash	settlement	accept	175	\N		\N	0
169	MANTRA-176-1781011411717440075	bank_transfer	settlement	accept	176	\N		\N	0
170	MANTRA-177-1781011411723517571	qris	settlement	accept	177	\N		\N	0
171	MANTRA-178-1781011411731255247	bank_transfer	cancel	accept	178	\N		\N	0
172		cash	settlement	accept	179	\N		\N	0
173	MANTRA-180-1781011411745011809	bank_transfer	settlement	accept	180	\N		\N	0
174	MANTRA-181-1781011411751954978	gopay	settlement	accept	181	\N		\N	0
175		cash	cancel	accept	182	\N		\N	0
176	MANTRA-183-1781011411765344509	gopay	settlement	accept	183	\N		\N	0
177	MANTRA-184-1781011411773255213	gopay	settlement	accept	184	\N		\N	0
178		cash	settlement	accept	185	\N		\N	0
179	MANTRA-186-1781011411787968214	gopay	settlement	accept	186	\N		\N	0
180	MANTRA-187-1781011411795288418	gopay	pending	accept	187	\N		\N	0
181	MANTRA-1-1781011411801530472	bank_transfer	settlement	accept	1	\N		\N	0
182	MANTRA-23-1781011411807965174	bank_transfer	settlement	accept	23	\N		\N	0
183	MANTRA-56-1781011411815205151	gopay	pending	accept	56	\N		\N	0
184	MANTRA-64-1781011411822269346	qris	settlement	accept	64	\N		\N	0
185	MANTRA-97-1781011411829002990	qris	cancel	accept	97	\N		\N	0
186	MANTRA-126-1781011411836093619	gopay	settlement	accept	126	\N		\N	0
187	MANTRA-158-1781011411842308657	gopay	pending	accept	158	\N		\N	0
188	MANTRA-188-1781011411849069717	bank_transfer	settlement	accept	188	\N		\N	0
189	MANTRA-189-1781011411856402008	qris	settlement	accept	189	\N		\N	0
190	MANTRA-190-1781011411862533050	gopay	pending	accept	190	\N		\N	0
191	MANTRA-191-1781011411868966401	bank_transfer	settlement	accept	191	\N		\N	0
192	MANTRA-192-1781011411875776409	qris	settlement	accept	192	\N		\N	0
193	MANTRA-193-1781011411881773948	gopay	settlement	accept	193	\N		\N	0
194		cash	settlement	accept	194	\N		\N	0
195	MANTRA-195-1781011411894534195	gopay	settlement	accept	195	\N		\N	0
196		cash	settlement	accept	196	\N		\N	0
197		cash	cancel	accept	197	\N		\N	0
198	MANTRA-198-1781011411918884991	bank_transfer	pending	accept	198	\N		\N	0
199	MANTRA-199-1781011411925144406	qris	settlement	accept	199	\N		\N	0
200		cash	settlement	accept	200	\N		\N	0
201	MANTRA-201-1781011411941596193	qris	pending	accept	201	\N		\N	0
202	MANTRA-202-1781011411947510933	qris	pending	accept	202	\N		\N	0
203		cash	settlement	accept	203	\N		\N	0
204	MANTRA-204-1781011411961277679	bank_transfer	pending	accept	204	\N		\N	0
205	MANTRA-205-1781011411968015473	bank_transfer	settlement	accept	205	\N		\N	0
206	MANTRA-206-1781011411975866641	qris	settlement	accept	206	\N		\N	0
207	MANTRA-207-1781011411982489607	gopay	settlement	accept	207	\N		\N	0
208		cash	settlement	accept	208	\N		\N	0
209	MANTRA-209-1781011411995933655	bank_transfer	settlement	accept	209	\N		\N	0
210		cash	settlement	accept	210	\N		\N	0
211	MANTRA-211-1781011412008382758	qris	settlement	accept	211	\N		\N	0
212	MANTRA-212-1781011412015608133	qris	settlement	accept	212	\N		\N	0
213	MANTRA-213-1781011412022007240	bank_transfer	settlement	accept	213	\N		\N	0
214	MANTRA-214-1781011412028031053	bank_transfer	settlement	accept	214	\N		\N	0
215	MANTRA-215-1781011412035188441	gopay	settlement	accept	215	\N		\N	0
216		cash	settlement	accept	216	\N		\N	0
217	MANTRA-217-1781011412047297952	gopay	settlement	accept	217	\N		\N	0
218		cash	settlement	accept	218	\N		\N	0
219	MANTRA-219-1781011412059885546	gopay	settlement	accept	219	\N		\N	0
220	MANTRA-220-1781011412066375848	gopay	settlement	accept	220	\N		\N	0
221	MANTRA-221-1781011412072845985	gopay	settlement	accept	221	\N		\N	0
222	MANTRA-222-1781011412078991468	qris	pending	accept	222	\N		\N	0
223	MANTRA-223-1781011412084983111	qris	settlement	accept	223	\N		\N	0
224		cash	settlement	accept	224	\N		\N	0
225	MANTRA-225-1781011412096589690	bank_transfer	settlement	accept	225	\N		\N	0
226		pending	settlement	accept	226	1		2026-06-09 20:23:47.025792+07	103102960
228		pending	pending	accept	228	\N		\N	0
229	MID-20260609-1	qris	pending	challenge	1	\N		\N	0
230		pending	pending	accept	229	\N		\N	0
231	MID-20260609-1	qris	pending	challenge	1	\N		\N	0
227	MID-20260609-227	qris	settlement	accept	227	2	trx-midtrans-001	2026-06-09 21:56:41.134667+07	24270556
232		cash	settlement	accept	1	\N		\N	0
233		pending	settlement	accept	230	1		2026-06-10 08:44:54.622927+07	122202236
\.


--
-- Data for Name: pengantaran; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pengantaran (id_pengantaran, public_id, waktu_pickup, waktu_sampai, last_latitude, last_longitude, foto_bukti_pengiriman, id_pesanan, id_kurir, id_status_pengantaran, id_ekspedisi) FROM stdin;
1	69eceae5-94fc-4d12-a219-1f16a3d2463d	2026-06-09 18:23:32.106956+07	2026-06-09 19:23:32.106957+07	-7.05141	110.438125	https://picsum.photos/400/400	2	1	4	1
2	cc655e05-2c30-47e4-911a-8f4ca6a9647e	2026-06-09 18:23:32.113579+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	3	1	2	1
3	e80adb21-a5f4-4169-acc0-6a152d533c59	2026-06-09 18:23:32.119465+07	2026-06-09 19:23:32.119465+07	-7.05141	110.438125	https://picsum.photos/400/400	4	1	4	1
4	ca5fb56c-48f1-46a9-9112-11ae3c5ff16a	2026-06-09 18:23:32.125525+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	6	1	2	1
5	97d0d088-d83c-4a65-bb4b-f73e1c321a89	2026-06-09 18:23:32.131575+07	2026-06-09 19:23:32.131576+07	-7.05141	110.438125	https://picsum.photos/400/400	9	1	4	1
6	815335e1-aa09-4312-b325-c246a1b1c2b0	2026-06-09 18:23:32.137329+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	11	1	2	1
7	0fccfb4c-f66d-4d1c-b91b-c050a4733f92	2026-06-09 18:23:32.143246+07	2026-06-09 19:23:32.143246+07	-7.05141	110.438125	https://picsum.photos/400/400	13	1	4	1
8	af7b41d5-5a08-49c4-b2d4-c64bdd11d387	2026-06-09 18:23:32.14913+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	17	1	2	1
9	38edcc2e-18c1-4243-8ada-13f4f142d96f	2026-06-09 18:23:32.155382+07	2026-06-09 19:23:32.155383+07	-7.05141	110.438125	https://picsum.photos/400/400	21	1	4	1
10	543b4c72-6e18-4af1-afbc-d72b918a3608	2026-06-09 18:23:32.160878+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	26	1	2	1
11	2f66d1cf-2517-423f-90b4-3079a0cdabd2	2026-06-09 18:23:32.167105+07	2026-06-09 19:23:32.167105+07	-7.05141	110.438125	https://picsum.photos/400/400	28	1	4	1
12	24ce6cba-7f1a-45dd-ab41-a53ade9ac8f6	2026-06-09 18:23:32.172831+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	31	1	2	1
13	e3419041-62ed-4250-a4d6-c96f79f7f2c9	2026-06-09 18:23:32.178584+07	2026-06-09 19:23:32.178584+07	-7.05141	110.438125	https://picsum.photos/400/400	33	1	4	1
14	eb2a693e-19cb-42f0-8890-5539f1e5e5a5	2026-06-09 18:23:32.185531+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	34	1	2	1
15	4e30e675-a966-4f07-9213-f155efac72f6	2026-06-09 18:23:32.191544+07	2026-06-09 19:23:32.191545+07	-7.05141	110.438125	https://picsum.photos/400/400	35	1	4	1
16	6b9ff01c-fc06-4374-8726-e2ffd20e79d1	2026-06-09 18:23:32.198247+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	36	1	2	1
17	83d227b3-7644-43b1-8edd-bd1939a44a1b	2026-06-09 18:23:32.204151+07	2026-06-09 19:23:32.204152+07	-7.05141	110.438125	https://picsum.photos/400/400	37	1	4	1
18	a2dcef74-a99f-46c2-8585-34fc8da18ded	2026-06-09 18:23:32.210165+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	38	1	2	1
19	49a3b8c2-616f-4dc0-add7-c3f6ce0bd698	2026-06-09 18:23:32.217327+07	2026-06-09 19:23:32.217327+07	-7.05141	110.438125	https://picsum.photos/400/400	42	1	4	1
20	541ca8f8-e95d-4d77-8af0-7ac624c5c667	2026-06-09 18:23:32.22261+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	45	1	2	1
21	82ccb61b-8425-4e83-8596-ec2eb78f595c	2026-06-09 18:23:32.234367+07	2026-06-09 19:23:32.234367+07	-7.05141	110.438125	https://picsum.photos/400/400	46	1	4	1
22	d11f9262-4074-4861-9dd6-41e1b231b99c	2026-06-09 18:23:32.240552+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	47	1	2	1
23	c304e5cd-4b4b-40e0-be63-7a5640935fe1	2026-06-09 18:23:32.246614+07	2026-06-09 19:23:32.246614+07	-7.05141	110.438125	https://picsum.photos/400/400	49	1	4	1
24	b4a44a0e-b187-4399-b225-bd6ba1a0e1f9	2026-06-09 18:23:32.252322+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	54	1	2	1
25	86e7208c-9765-4d63-8a4e-ae3bb798bb62	2026-06-09 18:23:32.258596+07	2026-06-09 19:23:32.258597+07	-7.05141	110.438125	https://picsum.photos/400/400	55	1	4	1
26	3ac960ab-7d54-4232-af2b-ba4a11a8f933	2026-06-09 18:23:32.264263+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	57	1	2	1
27	214f3cbc-a485-4537-a9b6-772ed5fce818	2026-06-09 18:23:32.269754+07	2026-06-09 19:23:32.269754+07	-7.05141	110.438125	https://picsum.photos/400/400	58	1	4	1
28	b49da994-3f68-4149-89d4-f0182470a4d1	2026-06-09 18:23:32.276573+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	59	1	2	1
29	f88c646e-0e25-4def-b9f6-412ae3a67272	2026-06-09 18:23:32.282885+07	2026-06-09 19:23:32.282885+07	-7.05141	110.438125	https://picsum.photos/400/400	65	1	4	1
30	d28efbb9-f31c-4390-a75b-c8f013ea5b43	2026-06-09 18:23:32.288882+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	66	1	2	1
31	a5f3b677-981a-4ddc-8ed5-58fe19336bbd	2026-06-09 18:23:32.294751+07	2026-06-09 19:23:32.294751+07	-7.05141	110.438125	https://picsum.photos/400/400	68	1	4	1
32	0965bce4-71b2-4efd-a0f5-531662c8fdcf	2026-06-09 18:23:32.300635+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	69	1	2	1
33	727909b2-d7f1-407d-8c82-48faf94ddc92	2026-06-09 18:23:32.306503+07	2026-06-09 19:23:32.306504+07	-7.05141	110.438125	https://picsum.photos/400/400	70	1	4	1
34	cd35324e-8ebe-4c0c-9512-8f72046ce117	2026-06-09 18:23:32.312269+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	71	1	2	1
35	6aa77205-ebe9-47da-b522-494bec8e88d1	2026-06-09 18:23:32.318846+07	2026-06-09 19:23:32.318846+07	-7.05141	110.438125	https://picsum.photos/400/400	72	1	4	1
36	d102ec4b-cf05-44b3-94f1-a4d4c15d24a4	2026-06-09 18:23:32.324726+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	74	1	2	1
37	055d64b9-f30d-4f8e-8b93-07861af02574	2026-06-09 18:23:32.330943+07	2026-06-09 19:23:32.330944+07	-7.05141	110.438125	https://picsum.photos/400/400	75	1	4	1
38	3f5939c5-7757-4b90-b1a0-159515bc8117	2026-06-09 18:23:32.33683+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	76	1	2	1
39	6c627e57-8c57-485e-adb2-09fa7eae56ee	2026-06-09 18:23:32.342661+07	2026-06-09 19:23:32.342661+07	-7.05141	110.438125	https://picsum.photos/400/400	78	1	4	1
40	62c136db-4cac-43f6-b435-02c8c6154805	2026-06-09 18:23:32.348894+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	79	1	2	1
41	3049d9ca-1be9-4545-86df-68ccf1141923	2026-06-09 18:23:32.354534+07	2026-06-09 19:23:32.354534+07	-7.05141	110.438125	https://picsum.photos/400/400	81	1	4	1
42	439da624-e430-4bec-908a-ffdc9b0670be	2026-06-09 18:23:32.3613+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	82	1	2	1
43	65ffb394-0549-4025-8640-6cb004e336c2	2026-06-09 18:23:32.36696+07	2026-06-09 19:23:32.36696+07	-7.05141	110.438125	https://picsum.photos/400/400	83	1	4	1
44	6783d029-eb29-4d5b-abe4-70a2810a14ab	2026-06-09 18:23:32.372792+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	86	1	2	1
45	214c8c33-38de-4e2e-b8e1-6dbf3dbadd48	2026-06-09 18:23:32.378667+07	2026-06-09 19:23:32.378667+07	-7.05141	110.438125	https://picsum.photos/400/400	87	1	4	1
46	35f59ba0-6aff-41b2-b70f-996d006b7ce8	2026-06-09 18:23:32.385455+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	88	1	2	1
47	2a53bc94-86a1-4b7b-867d-ca33e9155974	2026-06-09 18:23:32.392206+07	2026-06-09 19:23:32.392206+07	-7.05141	110.438125	https://picsum.photos/400/400	89	1	4	1
48	42dbbe7a-3f1b-4ca6-ab90-2f055d9cd395	2026-06-09 18:23:32.398033+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	90	1	2	1
49	526b94ae-ad73-429b-8558-2d74cfb00070	2026-06-09 18:23:32.40479+07	2026-06-09 19:23:32.404791+07	-7.05141	110.438125	https://picsum.photos/400/400	91	1	4	1
50	d1471df4-6311-4e56-8f2b-b5c2903f478a	2026-06-09 18:23:32.410522+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	98	1	2	1
51	bb98dbc7-ca37-4f76-a4a2-23c7c470469b	2026-06-09 18:23:32.416207+07	2026-06-09 19:23:32.416207+07	-7.05141	110.438125	https://picsum.photos/400/400	99	1	4	1
52	8c0b4491-6b3b-4cdc-b2d8-3e873b8f7e16	2026-06-09 18:23:32.424148+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	101	1	2	1
53	f3260663-1089-4832-8744-54e3876e9d06	2026-06-09 18:23:32.43127+07	2026-06-09 19:23:32.431271+07	-7.05141	110.438125	https://picsum.photos/400/400	102	1	4	1
54	39e99012-753d-44b1-b967-da4299779520	2026-06-09 18:23:32.436693+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	103	1	2	1
55	b3a0e29a-fab0-4435-8948-6e74a6b3512d	2026-06-09 18:23:32.442965+07	2026-06-09 19:23:32.442965+07	-7.05141	110.438125	https://picsum.photos/400/400	105	1	4	1
56	a40b12c1-53ba-4860-880e-9d7a28eee4be	2026-06-09 18:23:32.448763+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	107	1	2	1
57	23636da5-130f-41d3-b1fd-df64bcb5ac05	2026-06-09 18:23:32.455927+07	2026-06-09 19:23:32.455927+07	-7.05141	110.438125	https://picsum.photos/400/400	108	1	4	1
58	91c8431b-6e8d-4886-a296-d1a30ba05c72	2026-06-09 18:23:32.462697+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	111	1	2	1
59	d79f2602-8aa0-421d-8d2c-ed70aad48d1c	2026-06-09 18:23:32.468608+07	2026-06-09 19:23:32.468608+07	-7.05141	110.438125	https://picsum.photos/400/400	112	1	4	1
60	aa3a5820-ed49-450f-a596-d28f8b788091	2026-06-09 18:23:32.474688+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	113	1	2	1
61	af885010-a7f7-45f7-b855-d714d41fcdae	2026-06-09 18:23:32.481299+07	2026-06-09 19:23:32.4813+07	-7.05141	110.438125	https://picsum.photos/400/400	115	1	4	1
62	af0c248f-0d52-48dc-a538-d255d3fc215e	2026-06-09 18:23:32.487126+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	116	1	2	1
63	17259b3e-15d1-469c-89de-520d64f24f72	2026-06-09 18:23:32.492752+07	2026-06-09 19:23:32.492752+07	-7.05141	110.438125	https://picsum.photos/400/400	117	1	4	1
64	a1c62b4d-edaf-4010-a86b-74c5b29f09f3	2026-06-09 18:23:32.49857+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	118	1	2	1
65	02d32e1b-48f5-4e28-b8ff-92f91c74f324	2026-06-09 18:23:32.505389+07	2026-06-09 19:23:32.50539+07	-7.05141	110.438125	https://picsum.photos/400/400	119	1	4	1
66	a8614000-7257-46e9-8f91-85290f11d790	2026-06-09 18:23:32.511965+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	120	1	2	1
67	2463df47-bd6f-453a-bd63-2442402e7b72	2026-06-09 18:23:32.51773+07	2026-06-09 19:23:32.51773+07	-7.05141	110.438125	https://picsum.photos/400/400	122	1	4	1
68	81299e2c-6970-4e75-b2b2-e785e044299d	2026-06-09 18:23:32.523277+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	123	1	2	1
69	40efbea5-950f-4e11-a814-04a73563cbcb	2026-06-09 18:23:32.528918+07	2026-06-09 19:23:32.528919+07	-7.05141	110.438125	https://picsum.photos/400/400	125	1	4	1
70	f371db30-fda1-4646-9da9-68cd051bfdbd	2026-06-09 18:23:32.534782+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	127	1	2	1
71	446bbf5d-0b7c-4171-8f1c-77cdec48aa01	2026-06-09 18:23:32.541007+07	2026-06-09 19:23:32.541007+07	-7.05141	110.438125	https://picsum.photos/400/400	128	1	4	1
72	37180ef3-a54b-41f2-975f-8e98b5e14cc9	2026-06-09 18:23:32.546495+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	130	1	2	1
73	405daf79-cad0-4ede-b720-4aa1e916caad	2026-06-09 18:23:32.552585+07	2026-06-09 19:23:32.552585+07	-7.05141	110.438125	https://picsum.photos/400/400	131	1	4	1
74	d96e8115-2471-4293-9960-90f7d57a3150	2026-06-09 18:23:32.558061+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	132	1	2	1
75	38bfad9a-006f-47b2-bef7-2e2c4a624c04	2026-06-09 18:23:32.570372+07	2026-06-09 19:23:32.570372+07	-7.05141	110.438125	https://picsum.photos/400/400	133	1	4	1
76	37ba4f97-5f1b-4175-be9b-617efdcbc8e9	2026-06-09 18:23:32.587085+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	134	1	2	1
77	1feee8a5-33c7-44fa-bcc4-e2b289f7dd33	2026-06-09 18:23:32.608497+07	2026-06-09 19:23:32.608498+07	-7.05141	110.438125	https://picsum.photos/400/400	139	1	4	1
78	9da516c0-6c94-42d2-8204-f3ce8ec0385e	2026-06-09 18:23:32.613896+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	140	1	2	1
79	f53084c9-b6b5-4c94-b52f-9efb392fdb3d	2026-06-09 18:23:32.619853+07	2026-06-09 19:23:32.619866+07	-7.05141	110.438125	https://picsum.photos/400/400	142	1	4	1
80	55e202f7-04b6-499e-b945-3a93ec0c33bc	2026-06-09 18:23:32.625042+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	144	1	2	1
81	b781f18f-77f7-4023-8ef2-54ff156e8ef5	2026-06-09 18:23:32.631104+07	2026-06-09 19:23:32.631104+07	-7.05141	110.438125	https://picsum.photos/400/400	145	1	4	1
82	a399fbd1-f14c-4a36-a97d-268990dcd4b5	2026-06-09 18:23:32.636238+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	147	1	2	1
83	19f68922-6eec-4e2d-b98c-d74622017906	2026-06-09 18:23:32.641733+07	2026-06-09 19:23:32.641733+07	-7.05141	110.438125	https://picsum.photos/400/400	148	1	4	1
84	5b187bcf-d10e-467f-b3a5-d47b35a49377	2026-06-09 18:23:32.647759+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	149	1	2	1
85	64f33ebe-9ffa-46af-bc01-16093db6a5d4	2026-06-09 18:23:32.654075+07	2026-06-09 19:23:32.654075+07	-7.05141	110.438125	https://picsum.photos/400/400	150	1	4	1
86	a2bbcaa0-cf87-44bb-a9fd-9417a1756aca	2026-06-09 18:23:32.663493+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	151	1	2	1
87	64ab6717-3ad3-41b3-a156-9d5297acae4d	2026-06-09 18:23:32.669444+07	2026-06-09 19:23:32.669444+07	-7.05141	110.438125	https://picsum.photos/400/400	154	1	4	1
88	f91a0bcf-e838-41e4-9578-b7676ab91f4e	2026-06-09 18:23:32.674953+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	156	1	2	1
89	4914fc29-8230-4527-8018-c9ff64a9fbfc	2026-06-09 18:23:32.680889+07	2026-06-09 19:23:32.680889+07	-7.05141	110.438125	https://picsum.photos/400/400	159	1	4	1
90	b055b929-ef51-4019-873c-4ce8399608a5	2026-06-09 18:23:32.687332+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	160	1	2	1
91	e19f78a0-19b3-4d21-b256-ac7f35fdee2f	2026-06-09 18:23:32.692657+07	2026-06-09 19:23:32.692657+07	-7.05141	110.438125	https://picsum.photos/400/400	161	1	4	1
92	7b00b135-d8c5-41c3-afd3-fb894eeccd4e	2026-06-09 18:23:32.697917+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	164	1	2	1
93	f0f27213-f79e-4b96-ae5f-60ac5adbaac3	2026-06-09 18:23:32.704514+07	2026-06-09 19:23:32.704514+07	-7.05141	110.438125	https://picsum.photos/400/400	165	1	4	1
94	a360e341-fee9-405e-8eab-184e0b57b1b8	2026-06-09 18:23:32.710371+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	167	1	2	1
95	e65baa65-b9e4-4eb0-921d-f3b33f80bed1	2026-06-09 18:23:32.715871+07	2026-06-09 19:23:32.715871+07	-7.05141	110.438125	https://picsum.photos/400/400	170	1	4	1
96	df750efe-5186-415b-bef2-d604268cbbf3	2026-06-09 18:23:32.721186+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	171	1	2	1
97	a45880f1-6cfd-4df9-8848-707c392805ba	2026-06-09 18:23:32.72637+07	2026-06-09 19:23:32.726371+07	-7.05141	110.438125	https://picsum.photos/400/400	172	1	4	1
98	7bb3146f-7e82-4296-b0d2-59c1ec2cd8dd	2026-06-09 18:23:32.732202+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	174	1	2	1
99	6a246133-7b4a-4eaf-94e0-a129ff2f130a	2026-06-09 18:23:32.737819+07	2026-06-09 19:23:32.737819+07	-7.05141	110.438125	https://picsum.photos/400/400	176	1	4	1
100	04a802af-a845-4988-ac89-cce2654b8b20	2026-06-09 18:23:32.743688+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	177	1	2	1
101	2f8eced2-e70f-4438-b196-c4b78aa70c3e	2026-06-09 18:23:32.748973+07	2026-06-09 19:23:32.748973+07	-7.05141	110.438125	https://picsum.photos/400/400	178	1	4	1
102	3c57312b-774e-4305-a2f0-b3db5634c3b5	2026-06-09 18:23:32.755231+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	180	1	2	1
103	5b91d7ac-f389-475d-af8d-c30c9ef343aa	2026-06-09 18:23:32.760905+07	2026-06-09 19:23:32.760906+07	-7.05141	110.438125	https://picsum.photos/400/400	181	1	4	1
104	9fffac23-3089-407d-914c-3756b807096a	2026-06-09 18:23:32.766209+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	183	1	2	1
105	81a5f3b1-8476-48ea-9b59-eb4a86e38a8f	2026-06-09 18:23:32.771914+07	2026-06-09 19:23:32.771915+07	-7.05141	110.438125	https://picsum.photos/400/400	184	1	4	1
106	5f8bb4ae-73c4-41a9-b01b-73a60686ab35	2026-06-09 18:23:32.777496+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	186	1	2	1
107	786c810a-b0b0-47a5-abd7-5f59a44df777	2026-06-09 18:23:32.783172+07	2026-06-09 19:23:32.783172+07	-7.05141	110.438125	https://picsum.photos/400/400	187	1	4	1
108	0f3896aa-8f63-4b9d-8935-6d6419f8d2f0	2026-06-09 18:23:32.788592+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	1	1	2	1
109	e86f0f5c-5b00-4d86-a9ef-d7b18dd3f4c4	2026-06-09 18:23:32.794445+07	2026-06-09 19:23:32.794446+07	-7.05141	110.438125	https://picsum.photos/400/400	23	1	4	1
110	5751891c-1212-4bd3-9b17-7ae480f9ffc8	2026-06-09 18:23:32.800086+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	56	1	2	1
111	d865b25b-43e5-4435-99a8-e0739d3a3a8c	2026-06-09 18:23:32.805788+07	2026-06-09 19:23:32.805789+07	-7.05141	110.438125	https://picsum.photos/400/400	64	1	4	1
112	ab563495-5a76-4d39-aa9e-380d17c5532b	2026-06-09 18:23:32.810879+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	97	1	2	1
113	e4f13f70-a353-430c-acb8-e122cce08354	2026-06-09 18:23:32.816721+07	2026-06-09 19:23:32.816721+07	-7.05141	110.438125	https://picsum.photos/400/400	126	1	4	1
114	4dba912a-8d1d-4ea8-9690-8de29dc0d9f2	2026-06-09 18:23:32.822653+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	158	1	2	1
115	5153641a-065b-455d-baa1-258611250324	2026-06-09 18:23:32.827906+07	2026-06-09 19:23:32.827906+07	-7.05141	110.438125	https://picsum.photos/400/400	188	1	4	1
116	a0934114-0b34-403f-aa48-fb0068aaddae	2026-06-09 18:23:32.83375+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	189	1	2	1
117	e7583432-eee6-4ccd-8fbc-33c25ab61cbe	2026-06-09 18:23:32.838818+07	2026-06-09 19:23:32.838818+07	-7.05141	110.438125	https://picsum.photos/400/400	190	1	4	1
118	07c6d466-9277-419c-9012-d69247b089d8	2026-06-09 18:23:32.844571+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	191	1	2	1
119	5693fcc4-b8be-4fe3-ad83-3539e71edcc6	2026-06-09 18:23:32.849924+07	2026-06-09 19:23:32.849924+07	-7.05141	110.438125	https://picsum.photos/400/400	192	1	4	1
120	87af70b8-b09a-4af3-942f-d19eabf572c5	2026-06-09 18:23:32.856131+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	193	1	2	1
121	80ae983b-0822-4aa7-a81f-2cd51595456e	2026-06-09 18:23:32.861778+07	2026-06-09 19:23:32.861779+07	-7.05141	110.438125	https://picsum.photos/400/400	195	1	4	1
122	dfc516c7-ae4d-4f81-a81f-3b0b8214d49c	2026-06-09 18:23:32.867783+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	198	1	2	1
123	afbe23f5-179a-445f-b7ac-469b9a6e7e38	2026-06-09 18:23:32.872995+07	2026-06-09 19:23:32.872995+07	-7.05141	110.438125	https://picsum.photos/400/400	199	1	4	1
124	9313f56f-6aa9-4656-9288-2c3eb583e675	2026-06-09 18:23:32.87872+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	201	1	2	1
125	4cd0b578-6eb8-4be8-918e-4fe6ee038ac5	2026-06-09 18:23:32.884053+07	2026-06-09 19:23:32.884054+07	-7.05141	110.438125	https://picsum.photos/400/400	202	1	4	1
126	052180c7-5af4-4f81-b119-ebb79061e70d	2026-06-09 18:23:32.889552+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	204	1	2	1
127	001d8444-c497-478c-8647-6d27e934ff9f	2026-06-09 18:23:32.894895+07	2026-06-09 19:23:32.894896+07	-7.05141	110.438125	https://picsum.photos/400/400	205	1	4	1
128	80dfd788-07a9-4559-a960-24eb484dd08a	2026-06-09 18:23:32.90097+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	206	1	2	1
129	983c1ffa-7329-4cb6-9b60-818cde319fdd	2026-06-09 18:23:32.906452+07	2026-06-09 19:23:32.906452+07	-7.05141	110.438125	https://picsum.photos/400/400	207	1	4	1
130	3d769917-f498-4a9d-bb58-7d9c1dae39f0	2026-06-09 18:23:32.913446+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	209	1	2	1
131	fd2fd576-5aa8-4cba-91d1-5d22861fa524	2026-06-09 18:23:32.918718+07	2026-06-09 19:23:32.918718+07	-7.05141	110.438125	https://picsum.photos/400/400	211	1	4	1
132	bc2459f1-5534-4cde-adb0-3ca303636a6b	2026-06-09 18:23:32.924767+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	212	1	2	1
133	3ac8497d-6550-43af-ac04-ec260c987162	2026-06-09 18:23:32.930125+07	2026-06-09 19:23:32.930125+07	-7.05141	110.438125	https://picsum.photos/400/400	213	1	4	1
134	7d4e4632-9beb-4131-8bca-ea7cad890070	2026-06-09 18:23:32.936703+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	214	1	2	1
135	f2ba0c88-db30-4dd5-956c-64a05d97541f	2026-06-09 18:23:32.942558+07	2026-06-09 19:23:32.942558+07	-7.05141	110.438125	https://picsum.photos/400/400	215	1	4	1
136	bbedab38-b994-427b-863f-d8d5ebbf0488	2026-06-09 18:23:32.94823+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	217	1	2	1
137	8f9d2c9b-3200-450d-b17c-a74818ca5033	2026-06-09 18:23:32.953526+07	2026-06-09 19:23:32.953526+07	-7.05141	110.438125	https://picsum.photos/400/400	219	1	4	1
138	3b3f0dcd-cd20-4fee-a9ce-341f1a2fed16	2026-06-09 18:23:32.958765+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	220	1	2	1
139	616e5a9c-5d59-4d81-83f5-9c967a01ca11	2026-06-09 18:23:32.964021+07	2026-06-09 19:23:32.964021+07	-7.05141	110.438125	https://picsum.photos/400/400	221	1	4	1
140	7964e8c3-f28c-4ec9-9c3f-cfa52dee7fc6	2026-06-09 18:23:32.969254+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	222	1	2	1
141	ab07e899-512d-4477-880f-a41d3d96990b	2026-06-09 18:23:32.974738+07	2026-06-09 19:23:32.974738+07	-7.05141	110.438125	https://picsum.photos/400/400	223	1	4	1
142	1cfcf07e-5224-45d2-81d4-a326197caf94	2026-06-09 18:23:32.980833+07	\N	-7.05141	110.438125	https://picsum.photos/400/400	225	1	2	1
\.


--
-- Data for Name: pesanan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pesanan (id_pesanan, public_id, total_pembayaran, tanggal_pesanan, tipe_pesanan, status_pesanan, id_customer, id_kasir, id_alamat, id_ekspedisi, id_layanan_ekspedisi, ongkos_kirim, catatan) FROM stdin;
2	d3d4656f-37c2-448d-93b1-6f3b218261b9	2308952	2026-06-09 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
3	eba17ee1-09d7-4b87-94e4-3543cf7a91f4	15725517	2026-06-09 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
4	e7e9ffef-5aa8-4b29-aff2-1757e9450195	4505418	2026-06-09 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
5	af471bcc-06d2-4a0f-a3bd-d6d24340154a	3224149	2026-06-09 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
6	3ef6e033-57a9-4a24-b730-6c6c1e94fd1f	2830169	2026-06-08 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
7	1e4a1dfc-1a3f-442c-b6f2-cd3d21d058ee	30469145	2026-06-08 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
8	486062fd-b7dd-44ca-ae54-490d50718de3	7051590	2026-06-08 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
9	b76465fa-beb5-478c-9428-3a89c44cd683	54444839	2026-06-08 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
10	9e7e441a-2612-4d3a-8d38-30b506bb1057	38683402	2026-06-08 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
11	3b35b5c8-f015-42db-9f8e-4d392dd03457	998697	2026-06-07 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
12	cce54e2c-8be0-4e7b-8cb5-1c99998a8c51	3405250	2026-06-07 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
13	40881967-8c72-420c-98fb-269613dc8a6c	2676165	2026-06-07 20:23:22.72596+07	Online	Diproses	1	1	1	\N	\N	0	
14	e605d0dd-7642-4095-8fb8-83bb365acc65	40064373	2026-06-07 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
15	37541d9a-7d71-4797-b105-3308172ab63f	34145886	2026-06-07 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
16	f3f9b7a1-dcbf-4546-a61f-761212869431	88445926	2026-06-07 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
17	f70268cf-495d-455c-a34b-fdc408a258fc	3124435	2026-06-07 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
18	62a7adae-5397-449b-a76a-2f99daf219b5	2947304	2026-06-07 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
19	1c49f194-f9d8-493f-a783-889852d6cfed	26475253	2026-06-07 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
20	88146566-83e4-480f-8615-2dec5d53d837	28971272	2026-06-06 20:23:22.72596+07	Offline	Diproses	1	1	\N	\N	\N	0	
21	f69af942-7cf0-426c-8ecd-d5d5b466b6df	485793	2026-06-06 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
22	94b9ab6d-4523-42a6-b79c-43bb659ee3f9	2597481	2026-06-06 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
24	40a7f3c5-77a0-42ed-b320-14367842afaf	45100327	2026-06-06 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
25	c5fe961f-e199-4378-8a5a-20bb36f573e8	20247608	2026-06-06 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
26	a86593c6-6d09-4a56-bb67-55b1659c38a6	1263712	2026-06-06 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
27	fc46a5be-b979-4efd-b64d-5fc2b39438e8	13834166	2026-06-06 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
28	30392d3e-c012-4022-a62e-9bee31533a3c	88181977	2026-06-05 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
29	b61afa5b-f6b5-4a77-9bc6-dd3cd2bc5d7b	2123749	2026-06-05 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
30	0844b197-cda5-486e-8d0c-34783f498a80	333536	2026-06-05 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
31	4f903339-7382-4078-a16e-69839e0f20b7	4917153	2026-06-05 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
32	83981709-4f0d-49ef-b941-1c28abc44172	2690885	2026-06-05 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
33	daf64b9c-8a5e-48aa-87fa-2add77378012	1755069	2026-06-05 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
34	4c72fc25-b399-4df3-999d-925214ee217e	60059227	2026-06-05 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
35	fbb83dd4-5440-400d-9329-9d62adaead1d	2268684	2026-06-04 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
36	4a3df5b8-fe69-40ac-a243-24987db30a64	44714740	2026-06-04 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
37	8d57632b-c033-4343-9a47-6425961f2ce9	5708192	2026-06-04 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
38	6f4f13a0-4b33-4503-afd3-e31ce9119a10	16827882	2026-06-04 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
39	3bcf8a64-6e25-4de7-8364-f567e2e66d9b	3297416	2026-06-04 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
40	4f1bf47c-8cd5-4271-9170-b66c0b9fff63	2047147	2026-06-03 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
41	246a5290-0868-4202-bb2f-a43e503a1a68	977101	2026-06-03 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
42	59a03370-2ecb-41e6-858d-18b4de1d4daa	1271894	2026-06-03 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
43	fb295dfe-8715-4bb1-a00e-da47e4613d97	14136663	2026-06-03 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
44	7e56ed2c-5cc8-4ce3-8b63-11cd75801897	3377588	2026-06-03 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
45	c3581bff-92f7-46b7-806e-8ce5c17cc3f6	32444822	2026-06-03 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
46	a61685a7-5bc9-4fd2-8d9a-338675194b6e	1211505	2026-06-03 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
47	78374e8e-b9ae-415a-bf07-06b19f0ec7b9	76576517	2026-06-03 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
48	aa4cafe3-6f8f-454c-b241-38b17e8347f6	1787801	2026-06-03 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
49	19cbf38f-4030-4269-8aba-383899878691	17678722	2026-06-03 20:23:22.72596+07	Online	Dibatalkan	1	1	1	\N	\N	0	
50	f3e1610e-587c-467a-9142-a8e6adf2b666	39569275	2026-06-02 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
51	951b56e2-2e89-4bf4-b444-5e99e8dfcb58	16583914	2026-06-02 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
52	7f46f916-9614-467d-948f-f92bdfdfe5e0	60861574	2026-06-02 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
53	3b765dc5-ab55-4512-9ed6-ebb40ca51ce3	3133826	2026-06-02 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
54	2fd15f86-8ac8-4a5c-b001-8770149443ec	2118965	2026-06-02 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
55	cf0ad490-7498-4063-87d2-255dc7228736	16974186	2026-06-02 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
57	f3cba330-1880-499f-a451-f61c47f78a2a	2000855	2026-06-01 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
58	d0b864b6-f64c-40cd-a8d6-1baf889297e1	3905703	2026-06-01 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
59	f67a73e0-32a9-4283-9fac-ec7600ebfc8c	30918937	2026-06-01 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
60	999174a0-2a9c-493b-a80c-65afb91e9bd6	2096204	2026-06-01 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
61	eb15c655-1403-4957-97d5-859fcb8af4a1	32122430	2026-06-01 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
62	c59e27b1-4d12-4c06-95d5-35d95a479499	1556622	2026-06-01 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
63	b09b7e24-b0d4-425d-8172-ba1768d6b780	107470	2026-06-01 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
65	bff34928-2746-4a8a-a7c5-016647518009	1143927	2026-05-31 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
66	2dda8b34-9ce0-49ed-b41c-f5c7eab6e5be	9850035	2026-05-31 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
67	598f36ed-4066-43d3-bb52-3636a52d634a	514023	2026-05-31 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
68	95d59090-b1ce-4364-b23e-b2a0691f3f3a	12161164	2026-05-31 20:23:22.72596+07	Online	Diproses	1	1	2	\N	\N	0	
69	36aebe81-96af-46d6-9900-30f064a1180d	2691974	2026-05-31 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
70	524cce4a-7a51-44c9-984f-ae9a59f3ddd7	1530009	2026-05-31 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
71	6528a7c0-d73b-4528-ae2e-7eb0c4da8d1b	28276222	2026-05-30 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
72	3f0949d5-10a0-4bb0-b735-8b0587f0b636	39436610	2026-05-30 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
73	e806b57e-7ce5-426f-972e-9bfc6f48d95a	1207086	2026-05-30 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
74	78741517-5af1-43b4-a44a-d788fca61a42	81502311	2026-05-30 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
75	d3d9afab-2768-4863-9833-4c1752237ac2	123577280	2026-05-30 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
76	126e4b97-88d5-4745-9cd3-acf7b7214044	50834182	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
77	c29d8a58-4f14-4486-86b7-4f822e32f267	54197742	2026-05-29 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
78	9daab4bc-1b75-4d3e-8ff1-883bff9ff39d	39175194	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
79	74169e75-bc03-4dd2-aecd-f7f0e243c064	6675568	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
80	39fccd6a-7f27-4c2f-ac78-3122035b6da7	51502762	2026-05-29 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
81	88663707-ce31-447b-a8ea-011f9e939428	11900153	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
82	d9bb5d2a-fc3c-440e-b233-8ecf5e050af8	3297330	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
83	429ee5db-6df5-4896-ab66-57627e74f6a2	14348707	2026-05-29 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
84	e5bc9b59-982b-49fa-9ec5-d729640bc560	17747826	2026-05-29 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
85	31eac9f2-5e94-4f5c-8491-5951000402f9	6011854	2026-05-28 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
86	8f9bd6b5-d175-4991-a6f9-7c58778e7561	4012616	2026-05-28 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
87	b509d998-847f-401c-ae40-a98c4f473490	6592732	2026-05-28 20:23:22.72596+07	Online	Diproses	1	1	1	\N	\N	0	
88	c5c00d80-154d-49b3-acd3-64830c3415a5	16512017	2026-05-28 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
89	0ba109cb-2e52-4728-9161-c1841c057668	1127359	2026-05-28 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
90	3f482e89-f12d-475c-8d69-36bff4b21f28	2109047	2026-05-27 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
91	499a837f-4640-4e6b-8b7d-f46aded3e715	38975788	2026-05-27 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
92	7753217f-aec3-4eab-9813-cdc0e4e66d68	13630089	2026-05-27 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
93	0ae9106a-b12e-4d1d-b6e0-65b033afa6da	40530217	2026-05-27 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
94	bdfe5f9a-0626-4c04-acbf-1702d2f0e93d	1247810	2026-05-27 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
95	3399cfd3-de38-480a-9816-c7c251194a0a	3757447	2026-05-27 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
96	756f82d8-56d6-4c71-8b3b-fe7b1e870457	1929581	2026-05-27 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
98	713f9f63-2bc3-4751-82f6-671d7ce884d5	23225306	2026-05-27 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
99	71ca1189-a01e-48ba-b8a5-6ad791d1e037	84069503	2026-05-27 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
100	e5416a23-1f23-4199-916e-72741de4142c	300291	2026-05-26 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
101	91463e41-bdb1-4457-8a68-e105b9872b19	13844070	2026-05-26 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
102	44ff2ee6-1678-40f1-887f-96857dd01a5d	793475	2026-05-26 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
103	9b696406-d340-4ad4-a16b-d63524d4e531	2510808	2026-05-26 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
104	a329e7d0-1f11-43e1-914a-2084529c7c98	17144149	2026-05-26 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
105	1ebaf4aa-7a29-4ab6-a4f5-db4ee5be5464	31953721	2026-05-26 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
106	274b5d34-befa-434b-8818-661998e958e4	22520254	2026-05-26 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
107	dd96d727-7d52-4799-8944-a92173559c5a	5909112	2026-05-26 20:23:22.72596+07	Online	Dikemas	1	1	1	\N	\N	0	
108	1708e236-a5ac-4d46-8f2f-951af0f46f6a	15071850	2026-05-26 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
109	2600ef58-31c7-4801-b9ea-cc674926d7ec	1256300	2026-05-26 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
110	f31c6021-239e-4a1b-adb9-ffc2f787d1b6	67568872	2026-05-25 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
111	19a6a6ee-b7b8-412a-a2bb-408e7a6c5020	711680	2026-05-25 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
112	507e6e80-a208-4204-b889-3b1ada03c337	29806585	2026-05-25 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
113	9fb06860-9acb-41bf-b240-eb63fa473492	4981872	2026-05-25 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
114	cb4a54f4-62e7-43f7-9f4d-d3220ff94720	246888	2026-05-25 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
115	78db9a83-d683-423e-aef6-28bd0aae962f	4958059	2026-05-24 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
116	c4bfa7a7-5807-473c-8bbb-9c8f7810f655	2596866	2026-05-24 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
117	132f3bfd-b18f-48c9-93c1-f9eb27a361eb	27274808	2026-05-24 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
118	67a6d631-3dcc-4f7a-a8f6-6c44eafbac26	2237498	2026-05-24 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
119	95ad0145-03a2-4596-809e-79c6f5c5bab9	5408732	2026-05-24 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
120	b63a4fbc-2320-4d66-8536-291305f1777e	356475	2026-05-24 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
121	5f53666a-5a92-4305-b8b9-cd720ec9f3c4	2688327	2026-05-24 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
122	b1ca3767-065d-4ef0-8b34-a503a3dfa1a3	2148304	2026-05-23 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
123	0029e9be-59ac-47ec-8809-c973341c0492	44981645	2026-05-23 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
124	ab6c464a-5dfc-416c-91b6-9cc8fba25fe2	25605351	2026-05-23 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
125	bafff75f-e897-4ded-a828-b0e86a636fcb	1596980	2026-05-23 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
127	95f4b4f1-e04f-42c1-abef-1d3a3b96463b	21199142	2026-05-23 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
128	f6522696-fd88-493b-8048-b403f363750b	389427	2026-05-23 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
129	e0a5e819-a727-4aa3-a457-79f3c11e93dc	20912065	2026-05-23 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
130	14cf9df6-f712-44d5-bf04-9b0960401a7f	30230378	2026-05-23 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
131	31e8e9ca-917e-48c0-8b10-3e0b63269453	3174173	2026-05-23 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
132	bb8f3750-57c7-4298-9e38-58caf84e68d0	4209126	2026-05-22 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
133	38af2c85-1ac2-4b6a-81f8-a36fe4173272	2007499	2026-05-22 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
134	42e3024b-33ca-4815-9891-e793511cb5bd	10204678	2026-05-22 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
135	c07efd86-4b03-4eb1-82a5-5c9a5dfc7214	3856930	2026-05-22 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
136	4c429272-263f-4a92-9684-c59bc1d6d952	65494583	2026-05-22 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
137	d479a73f-a196-46b7-bfa7-5102e293d6ae	17404014	2026-05-22 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
138	2d0659f9-9da7-4e08-87d4-f0b9a15422dd	3301720	2026-05-22 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
139	1f8f3851-2bed-42dd-8866-000b52521f0b	3682112	2026-05-21 20:23:22.72596+07	Online	Diproses	1	1	1	\N	\N	0	
140	8b31670a-8ad5-42a2-80dd-eab14856e679	29270201	2026-05-21 20:23:22.72596+07	Online	Diproses	1	1	2	\N	\N	0	
141	369a16b8-e749-4386-bc71-3fba26632118	3228549	2026-05-21 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
142	ed6f052b-c001-4ab9-bb2e-3b40c15c5a67	5158082	2026-05-21 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
143	6fbc5b8d-fe6b-4a1a-9dbe-a4804d6adbcd	1242566	2026-05-21 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
144	157bbd42-d973-4a6c-8d27-b69abeee5b22	65675650	2026-05-21 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
145	645ec3cf-2102-4668-99bf-bd19f1dcd5fb	1158109	2026-05-21 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
146	aedee613-dd8e-4e14-b4da-ef29bd125782	7071014	2026-05-21 20:23:22.72596+07	Offline	Dikemas	1	1	\N	\N	\N	0	
147	87c655a4-c3db-4028-aad5-30aba5a3e476	23667271	2026-05-21 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
148	873f11cd-7542-4e3b-b31b-e4446a64d91a	1403446	2026-05-20 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
149	21593fbc-0705-4f45-a83c-b46fc209895c	2504983	2026-05-20 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
150	ab769643-f1e4-4f44-91e9-1dedd4dcb103	1994882	2026-05-20 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
151	ccdf520f-1e0c-4951-bfb5-e71d5d47e5cd	1758723	2026-05-20 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
152	1d2bcadc-775e-4dce-b97a-dd544524327d	2092352	2026-05-20 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
153	b917ccae-b6c5-4620-a384-72d319235488	3471951	2026-05-20 20:23:22.72596+07	Offline	Diproses	1	1	\N	\N	\N	0	
154	fae89fb5-1821-45d2-8c21-0d01333a55d7	6501236	2026-05-20 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
155	c0ef8f7c-1bac-40a5-a8b3-94fe10bfb999	67364594	2026-05-20 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
156	985b2b2c-f986-4b1c-9c6f-ba02b56e90ae	70658531	2026-05-20 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
157	e67de79c-09b5-40ff-adf7-1502398a36dc	53383860	2026-05-19 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
159	71f18a2a-0471-443f-bd92-6dbd4a93d4b5	17804615	2026-05-19 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
160	3e2a849c-e23d-4a81-a52b-38b0e2f3bf9f	3386163	2026-05-19 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
161	deae25ef-f929-430a-8527-1e4feb37c40d	17744641	2026-05-19 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
162	27bd658c-0718-4ad1-be29-92809c25b551	79102180	2026-05-18 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
163	7426dc96-0dc8-4acc-9a00-3413d577188d	739370	2026-05-18 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
164	ee9db3e2-0214-48fa-92d0-ec2faa007f41	14207981	2026-05-18 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
165	20058866-7a6b-4302-b655-0b0ab020a599	1449422	2026-05-18 20:23:22.72596+07	Online	Diproses	1	1	1	\N	\N	0	
166	e5a53714-356b-48df-8768-a071eba6b303	5670183	2026-05-18 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
167	cb6f0a46-1c85-4031-97a9-dc5a97d62c54	1159600	2026-05-18 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
168	4945651b-2a4d-4efd-a2fd-c8539d16ec90	18123471	2026-05-18 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
169	6532ad1f-0216-4b3b-b241-e9ce30067dec	2704725	2026-05-18 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
170	68a7b010-bec9-46b4-9e15-ffdafe9f51c3	64584110	2026-05-17 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
171	fe1b9858-4475-4f5c-ab58-db65bd33f55c	5246729	2026-05-17 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
172	67b5ba94-c169-4ebb-ac93-680f0d19c2fa	12303054	2026-05-17 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
173	fed262e4-b3da-4985-8756-760a4fc55953	30323830	2026-05-17 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
174	97ab11a5-ff1b-4354-bf84-2554700bf68a	51875560	2026-05-17 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
175	bd969531-04fe-4bae-b07a-d20cbcb1965d	1857722	2026-05-17 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
176	917b61eb-ea8e-46b5-847c-600565c1f67d	1994467	2026-05-17 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
177	d670600b-3f56-4556-97df-bd4002c9d4f6	4301955	2026-05-16 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
178	cd3fd27a-ce5e-468c-9127-a0372de12be2	16604064	2026-05-16 20:23:22.72596+07	Online	Dibatalkan	1	1	2	\N	\N	0	
179	7672c385-1494-4ad8-8454-938d2f619e1e	16946262	2026-05-16 20:23:22.72596+07	Offline	Dikirim	1	1	\N	\N	\N	0	
180	1299d59d-0d30-4d0b-859f-e5bf3e1c2b31	2516981	2026-05-16 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
181	761fe7b1-e4bf-4db7-80f6-3dbb2682e59b	2392776	2026-05-16 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
182	00184518-935c-49ad-8bed-fd248864313e	1752905	2026-05-16 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
183	029b315d-b646-49fb-9d98-d2978bbed3f6	61099707	2026-05-16 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
184	58f6c597-bcf0-4b26-a138-f441c89eff2a	27254297	2026-05-16 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
185	d7145165-b09d-40b4-9a79-469fc06d35d3	3635370	2026-05-16 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
186	8a0433aa-80bc-471e-a8dc-193890c70ba8	965790	2026-05-15 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
187	9370705f-caf1-4186-b8af-1f5906a31a3b	1744792	2026-05-15 20:23:22.72596+07	Online	Diproses	1	1	1	\N	\N	0	
23	79f519fd-0d56-4e2b-806e-013530521e66	1654084	2026-06-06 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
56	1acf8f7a-19ba-4e9f-b3c1-cb28a809f6b9	20968818	2026-06-02 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
64	eff9e2c6-804e-485d-962f-1d396df854b0	2065380	2026-06-01 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
97	72ad011a-f177-406a-8b60-1a856390b024	1552858	2026-05-27 20:23:22.72596+07	Online	Dibatalkan	1	1	1	\N	\N	0	
126	f3adbace-1811-4264-aa4f-128cf721035a	982090	2026-05-23 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
158	d55e069f-c73d-4b13-aa7c-04b88aa863a7	613773	2026-05-19 20:23:22.72596+07	Online	Diproses	1	1	2	\N	\N	0	
188	9bb6f26c-919a-4970-953a-4ab8c06cf82b	2424049	2026-05-15 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
189	596a6a9f-4c9f-4e83-9507-5b8b26b7da61	60059578	2026-05-15 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
190	fbf649cd-67cd-4d4b-b65a-7bf3c4c71729	1493725	2026-05-15 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
191	7711f32f-adff-4d6b-89aa-540bf076c6bf	1773104	2026-05-15 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
192	890050fc-7bd3-4e75-962d-054a248326db	5918399	2026-05-14 20:23:22.72596+07	Online	Dikirim	1	1	2	\N	\N	0	
193	2c2aade3-c9d4-4a01-9f3e-56de9434236a	62705051	2026-05-14 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
194	204e7928-a9e9-409f-8dfa-f9e7e57e6eab	2137563	2026-05-14 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
195	86d5742d-78df-4444-9833-99ff39ee7d57	1660551	2026-05-14 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
196	89837cff-084e-4833-b741-c31e358773dc	11055654	2026-05-14 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
197	bf0661ab-022a-4f2f-a83a-c8a6dba2ec6d	70650196	2026-05-14 20:23:22.72596+07	Offline	Dibatalkan	1	1	\N	\N	\N	0	
198	1f5f1401-1f2a-46f2-b5ae-ffcd9cceef5c	639328	2026-05-14 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
199	cc569696-3fab-4365-a594-28628dd2d0a8	56587426	2026-05-14 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
200	2b9d843b-4fc3-423c-9c71-066dad755a06	534269	2026-05-14 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
201	4f8aefbe-5b4e-44ea-97f1-30ba6352601f	979008	2026-05-14 20:23:22.72596+07	Online	Dikemas	1	1	1	\N	\N	0	
202	19f128ab-dc43-4a22-b67e-4ef4c2e2f1b5	2594809	2026-05-13 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
203	3f2a7740-15b4-4292-bc4f-a30703a4c150	47778866	2026-05-13 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
204	125f84e1-1696-49e0-a988-5076346f43dd	890475	2026-05-13 20:23:22.72596+07	Online	Dikemas	1	1	2	\N	\N	0	
205	c8084759-ee59-4996-b0b0-5fd1e32bd345	1868744	2026-05-13 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
206	b097f6a3-5224-46ba-a176-0cf06839c6d0	2203835	2026-05-13 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
207	8f0e6b9d-56f9-47f2-b921-5ebc4baf1f8a	75452114	2026-05-13 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
208	60d86f46-400d-4a71-8588-b3c86b90e3ae	1148362	2026-05-13 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
209	18fe27c9-ca6b-4323-b2d0-5bf3e27ba7fc	2207425	2026-05-13 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
210	0d37a282-c34c-4e09-ae5b-6980ea823809	1904570	2026-05-12 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
211	aaa2da6b-c744-4ea7-b996-4b3a8c4cdbef	2207931	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
212	0b8ba012-8e8c-4169-a2d9-428b7cc89059	11698055	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
213	375a96bd-1021-4f99-951d-39ffbfb4add5	53493290	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
214	6fc460ea-d0ae-407d-a76f-cb3abeaf8a30	1117934	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
215	e0fb7afb-e782-4755-9034-be027dfd2b7e	1756258	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
216	48a16e26-dea3-49f3-892b-4fcd4eeb1657	75905225	2026-05-12 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
217	56db94a7-c466-4547-a122-774bc1f999a0	46768191	2026-05-12 20:23:22.72596+07	Online	Dikirim	1	1	1	\N	\N	0	
218	91078b61-bef7-4014-9cdc-ebec885ce312	16703197	2026-05-12 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
219	2185b832-2c5a-403f-873d-2168b5aedd5d	3219123	2026-05-12 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
220	b565ca23-235d-47df-8435-ac465ac2229f	323549	2026-05-11 20:23:22.72596+07	Online	Selesai	1	1	2	\N	\N	0	
221	ac643c8c-1e44-4c22-94c9-fbb87eaa06d3	70566582	2026-05-11 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
222	359b9e4f-b2af-48c1-a098-a222d8e13bd0	1837563	2026-05-11 20:23:22.72596+07	Online	Diproses	1	1	2	\N	\N	0	
223	70705cc3-a491-4e17-b75f-254a6d916ae0	917976	2026-05-11 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
224	39d6bb6c-225b-40d2-830a-1bed5829c91b	884589	2026-05-11 20:23:22.72596+07	Offline	Selesai	1	1	\N	\N	\N	0	
225	4ae00e40-ee72-4ddc-9265-09ba47229141	15691290	2026-05-11 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
226	c07a516d-2a0f-4ebd-afcd-a2ef9afd2962	103102960	2026-06-09 20:23:47.010075+07	Online	Selesai	1	\N	1	1	1	0	
228	b6e252eb-94b2-4672-954f-16e8d56f9de2	24270556	2026-06-09 20:25:27.015329+07	Online	Diproses	1	\N	\N	\N	\N	0	
229	216b5e3e-ef66-45ee-93b0-ea73e410154f	103102960	2026-06-09 21:14:32.828753+07	Online	Diproses	1	\N	\N	\N	\N	0	
227	891958ba-217f-472d-ab8d-a34760be2461	24270556	2026-06-09 20:24:06.215364+07	Online	Selesai	1	\N	1	1	1	0	
1	2d444438-dc02-4ab1-9320-9cc61ed16374	3864032	2026-06-09 20:23:22.72596+07	Online	Selesai	1	1	1	\N	\N	0	
230	4c199f68-e023-41f4-9ea6-e901fc4a808d	122202236	2026-06-10 08:44:54.608194+07	Online	Selesai	1	\N	1	1	1	15000	Test E2E checkout
\.


--
-- Data for Name: refresh_token; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.refresh_token (id_token, token, expires_at, created_at, revoked_at, id_user) FROM stdin;
1	6091610ffa8ea6214666944b3c0ba528ecefc99fdf56f423379c45c77ee0129e	2026-06-16 20:23:34.24398+07	2026-06-09 20:23:34.24398+07	\N	4
2	109c6d3a506c33c17775c60e89789be28fe1f7afbb977400428abad112bd008d	2026-06-16 20:23:40.889715+07	2026-06-09 20:23:40.889715+07	\N	4
3	cb0b6d377288974deebbec8e2a5a48f2ee9ba6d9450d47b0ea9013a2e1f5bb94	2026-06-16 20:24:52.743246+07	2026-06-09 20:24:52.743246+07	2026-06-09 20:25:06.65348+07	1
4	807898438af1bb988aa38014a853801d608b1bff751c4cad8988980cbeb9cacc	2026-06-16 20:25:17.328362+07	2026-06-09 20:25:17.328362+07	2026-06-09 20:25:21.196793+07	4
5	ecc208d32b0164a0d3fa2e736e076ee4bc8332e469fe0e5e8391e325af28bb3d	2026-06-16 20:25:30.756278+07	2026-06-09 20:25:30.756279+07	2026-06-09 20:25:31.370825+07	2
6	46b373d520ea45961a943bf02f653426a9015b54facbaf1b765650a9a4f09a38	2026-06-16 20:25:31.590173+07	2026-06-09 20:25:31.590173+07	\N	3
7	fec5fa50d4ab59df9d8dbdc66922a28095097e2377fda47ff242670e006f91a2	2026-06-16 21:14:05.057303+07	2026-06-09 21:14:05.057303+07	2026-06-09 21:14:19.795261+07	1
8	28ccaba307ff2cf622d4151ccce75cf877bd1cdada1daa805574bb6acaf3bbe5	2026-06-16 21:14:24.742006+07	2026-06-09 21:14:24.742006+07	2026-06-09 21:14:26.793679+07	4
9	47d9d7137aef78a6d6c9d66862096dd28e12b7d35220ef924290101e9f6a8d0d	2026-06-16 21:14:35.899123+07	2026-06-09 21:14:35.899123+07	2026-06-09 21:14:36.519632+07	2
10	c93900f6fb312b07afc617068fe3d7caf43805273304c8623e6c5eab815242ae	2026-06-16 21:14:43.388996+07	2026-06-09 21:14:43.388996+07	\N	3
11	3d66a875aa2fd4620535ffad601d684de80c513c2500a64e46bdb54168a177a9	2026-06-16 21:29:27.949978+07	2026-06-09 21:29:27.949978+07	\N	3
12	b56ff89db823ee0c805c2222e1489abe1e19fd4c96e25b79f61ce815fed59058	2026-06-16 21:29:35.093573+07	2026-06-09 21:29:35.093574+07	\N	2
13	02526bb77eede74a6d0268fa85b248b9280669ddd76ac0342fd541b62012ee60	2026-06-16 21:33:53.259113+07	2026-06-09 21:33:53.259113+07	\N	2
14	8067531c0641ad531c4af4ca31857a4b698034e71675301cdfda9e59420fdbd6	2026-06-16 21:55:42.172651+07	2026-06-09 21:55:42.172652+07	\N	3
15	97a3b52d34720b1f0cc78b73d2a997c6bcc56a0dbaba331a3a80668ee546cded	2026-06-16 22:07:55.700207+07	2026-06-09 22:07:55.700207+07	\N	4
16	8af68c49486332f5966e44f4108c443345b88151ca7443a97400131530836e1a	2026-06-16 23:21:30.20689+07	2026-06-09 23:21:30.206891+07	\N	2
17	0c4369d92edfbc6197d64fb93bd5bf735c01df7e47606905c4f84f511c6274f3	2026-06-16 23:43:27.621181+07	2026-06-09 23:43:27.621181+07	\N	2
18	a6585213fb54e1fc10e6ce9c9ff728764536486424c3da1f9b9e2c79504feabf	2026-06-17 07:20:06.008519+07	2026-06-10 07:20:06.00852+07	2026-06-10 07:20:52.309846+07	1
19	15382adc126803c8e634828ba41229c7779ac6edbd91cca8fbb8838b73e3f32f	2026-06-17 07:32:08.399335+07	2026-06-10 07:32:08.399335+07	\N	1
20	8736417effcfe4a7773acb0704f1b632d795c8e3f4fc556e9ad3aaab13c68197	2026-06-17 07:37:36.463398+07	2026-06-10 07:37:36.463398+07	\N	4
21	957434a0117d516be158194469c163ef9a73967a333158f0ffe920478e4eed22	2026-06-17 08:08:28.062129+07	2026-06-10 08:08:28.062129+07	\N	2
22	207a36b9c4b6d4c09285c1cfa4c597bd74d7d94c947f244c56f00badd4596f80	2026-06-17 08:15:43.787744+07	2026-06-10 08:15:43.787744+07	2026-06-10 08:16:21.188014+07	1
23	8e6a8b10a69adb27e5686c635874c4c8374e7570bd0131ab91aceb5bf7fc7726	2026-06-17 08:16:29.352758+07	2026-06-10 08:16:29.352758+07	2026-06-10 08:20:20.700446+07	1
24	75e28a9944a1f8b91701b1a49a3069177e6f07c4fe092f3ce0592d6b3424cf75	2026-06-17 08:42:45.012558+07	2026-06-10 08:42:45.012558+07	\N	4
25	d64170affce8bfe6a860aa75282489fb064f89555027d8d5e7ad1d3a4c3586e6	2026-06-17 08:43:09.881592+07	2026-06-10 08:43:09.881592+07	\N	4
26	c883c32d18cd8689c2714d0c589f7af981c88cd585d52cdf5bc6497b3f136045	2026-06-17 08:43:42.71468+07	2026-06-10 08:43:42.71468+07	\N	4
\.


--
-- Data for Name: role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.role (id_role, nama_role) FROM stdin;
1	Customer
2	Admin
3	Kasir
4	Kurir
\.


--
-- Data for Name: satuan; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.satuan (id_satuan, nama_satuan) FROM stdin;
1	Pcs
2	Kg
3	Gram
4	Cup
\.


--
-- Data for Name: spesifikasi; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.spesifikasi (id_spesifikasi, nama_spesifikasi) FROM stdin;
1	Warna
2	Ukuran
\.


--
-- Data for Name: spesifikasi_barang; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.spesifikasi_barang (id_spesifikasi_barang, jumlah, harga_barang, berat_barang, id_barang, id_detail_spesifikasi) FROM stdin;
3	42	14097752	0	2	3
4	33	5588932	0	2	4
5	44	10002883	0	3	4
6	5	3229255	0	3	5
7	48	626578	0	4	5
8	9	263949	0	4	6
9	39	775529	0	5	6
10	49	455250	0	5	7
11	1	776868	0	6	7
12	42	775768	0	6	8
13	48	40359	0	7	8
14	27	105736	0	7	9
15	9	52841	0	8	9
16	1	162448	0	8	10
17	44	196418	0	9	10
18	41	193009	0	9	11
19	32	97462	0	10	11
20	49	67111	0	10	12
21	5	30543	0	11	12
22	7	21363	0	11	13
23	25	38911	0	12	13
24	43	54842	0	12	1
25	18	455963	0	13	1
26	2	241842	0	13	2
27	17	423793	0	14	2
28	41	112185	0	14	3
29	7	240763	0	15	3
30	26	172895	0	15	4
1	2	12861980	0	1	2
2	48	14816845	0	1	3
\.


--
-- Data for Name: status_pengantaran; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.status_pengantaran (id_status_pengantaran, nama_status) FROM stdin;
1	Menunggu Pickup
2	Dalam Perjalanan
3	Tiba di Tujuan
4	Selesai
5	Gagal Antar
\.


--
-- Data for Name: stok_opname; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stok_opname (id_stok_opname, harga_beli, status, jumlah_stok, keterangan, tanggal, id_spesifikasi_barang) FROM stdin;
1	12856997	t	57	Stok keluar penjualan	2026-05-25 20:23:22.140562+07	1
2	12858923	f	17	Stok opname bulanan	2026-06-03 20:23:22.147031+07	1
3	12856225	t	23	Stok masuk dari supplier	2026-05-31 20:23:22.153187+07	1
4	14813033	t	99	Stok masuk dari supplier	2026-05-10 20:23:22.159088+07	2
5	14809535	f	14	Retur barang	2026-05-26 20:23:22.16509+07	2
6	14808847	t	45	Stok opname bulanan	2026-05-11 20:23:22.171339+07	2
7	14087761	t	41	Stok opname bulanan	2026-05-31 20:23:22.177805+07	3
8	14091034	f	23	Retur barang	2026-05-10 20:23:22.183449+07	3
9	14091496	t	43	Retur barang	2026-05-18 20:23:22.190195+07	3
10	5575174	t	23	Stok opname bulanan	2026-05-29 20:23:22.196881+07	4
11	5574949	f	34	Stok opname bulanan	2026-05-30 20:23:22.201921+07	4
12	5578668	t	78	Stok masuk dari supplier	2026-06-02 20:23:22.208374+07	4
13	9989889	t	27	Stok opname bulanan	2026-05-11 20:23:22.214316+07	5
14	9996884	f	19	Retur barang	2026-05-17 20:23:22.220255+07	5
15	9993767	t	82	Stok keluar penjualan	2026-05-30 20:23:22.22623+07	5
16	3214279	t	93	Stok keluar penjualan	2026-05-28 20:23:22.232257+07	6
17	3223412	f	79	Stok opname bulanan	2026-05-13 20:23:22.238197+07	6
18	3220602	t	92	Retur barang	2026-06-07 20:23:22.244487+07	6
19	613733	t	43	Stok opname bulanan	2026-05-26 20:23:22.250546+07	7
20	625113	f	81	Retur barang	2026-05-18 20:23:22.256359+07	7
21	622567	t	26	Stok opname bulanan	2026-05-23 20:23:22.262491+07	7
22	256653	t	17	Stok masuk dari supplier	2026-05-18 20:23:22.268654+07	8
23	255845	f	77	Retur barang	2026-05-27 20:23:22.274175+07	8
24	257462	t	41	Stok masuk dari supplier	2026-05-25 20:23:22.281268+07	8
25	764289	t	69	Stok keluar penjualan	2026-05-22 20:23:22.287206+07	9
26	761476	f	53	Stok opname bulanan	2026-05-12 20:23:22.292346+07	9
27	768830	t	15	Stok masuk dari supplier	2026-05-22 20:23:22.298189+07	9
28	444129	t	26	Retur barang	2026-06-01 20:23:22.303833+07	10
29	440516	f	78	Stok opname bulanan	2026-05-31 20:23:22.309166+07	10
30	449225	t	54	Retur barang	2026-05-14 20:23:22.315278+07	10
31	763837	t	95	Retur barang	2026-06-01 20:23:22.321188+07	11
32	768236	f	76	Stok keluar penjualan	2026-05-14 20:23:22.327587+07	11
33	765775	t	30	Stok masuk dari supplier	2026-06-08 20:23:22.33452+07	11
34	766616	t	18	Retur barang	2026-05-25 20:23:22.340206+07	12
35	765003	f	82	Stok keluar penjualan	2026-05-22 20:23:22.345878+07	12
36	768526	t	87	Stok keluar penjualan	2026-05-15 20:23:22.352209+07	12
37	26764	t	53	Retur barang	2026-05-28 20:23:22.35791+07	13
38	33575	f	13	Stok keluar penjualan	2026-06-02 20:23:22.363244+07	13
39	30094	t	53	Stok keluar penjualan	2026-05-20 20:23:22.369778+07	13
40	92125	t	46	Stok masuk dari supplier	2026-05-11 20:23:22.376284+07	14
41	99917	f	15	Stok keluar penjualan	2026-05-24 20:23:22.382074+07	14
42	99238	t	60	Stok masuk dari supplier	2026-05-23 20:23:22.405339+07	14
43	47657	t	19	Stok keluar penjualan	2026-06-02 20:23:22.429602+07	15
44	51586	f	55	Stok masuk dari supplier	2026-05-26 20:23:22.435717+07	15
45	39326	t	81	Stok opname bulanan	2026-06-04 20:23:22.442455+07	15
46	157342	t	98	Stok opname bulanan	2026-05-23 20:23:22.44879+07	16
47	149399	f	60	Stok masuk dari supplier	2026-05-26 20:23:22.455344+07	16
48	148396	t	75	Stok masuk dari supplier	2026-05-20 20:23:22.462357+07	16
49	186041	t	45	Retur barang	2026-05-23 20:23:22.46959+07	17
50	183898	f	80	Stok keluar penjualan	2026-05-20 20:23:22.476096+07	17
51	186168	t	93	Stok keluar penjualan	2026-05-21 20:23:22.483333+07	17
52	181571	t	65	Stok keluar penjualan	2026-05-18 20:23:22.490428+07	18
53	184084	f	95	Retur barang	2026-05-11 20:23:22.496952+07	18
54	183658	t	31	Stok keluar penjualan	2026-06-04 20:23:22.503727+07	18
55	92158	t	24	Stok keluar penjualan	2026-05-22 20:23:22.510271+07	19
56	87572	f	74	Stok keluar penjualan	2026-06-07 20:23:22.515931+07	19
57	89328	t	88	Retur barang	2026-05-11 20:23:22.522239+07	19
58	57931	t	65	Stok keluar penjualan	2026-05-11 20:23:22.527985+07	20
59	65464	f	48	Retur barang	2026-06-03 20:23:22.533418+07	20
60	54876	t	80	Retur barang	2026-06-01 20:23:22.539704+07	20
61	24183	t	85	Retur barang	2026-05-11 20:23:22.546212+07	21
62	20885	f	30	Retur barang	2026-06-02 20:23:22.552728+07	21
63	22220	t	96	Stok masuk dari supplier	2026-05-17 20:23:22.558829+07	21
64	9711	t	54	Stok opname bulanan	2026-05-28 20:23:22.564634+07	22
65	8142	f	23	Stok masuk dari supplier	2026-06-08 20:23:22.570063+07	22
66	15313	t	31	Stok keluar penjualan	2026-05-15 20:23:22.575821+07	22
67	31173	t	59	Stok keluar penjualan	2026-05-25 20:23:22.581503+07	23
68	36955	f	66	Retur barang	2026-05-15 20:23:22.587141+07	23
69	30929	t	46	Stok keluar penjualan	2026-05-15 20:23:22.59345+07	23
70	51527	t	83	Stok masuk dari supplier	2026-05-28 20:23:22.600032+07	24
71	42914	f	49	Retur barang	2026-05-24 20:23:22.606041+07	24
72	49468	t	44	Stok keluar penjualan	2026-05-19 20:23:22.612313+07	24
73	446001	t	58	Retur barang	2026-05-30 20:23:22.618343+07	25
74	441644	f	100	Stok masuk dari supplier	2026-05-24 20:23:22.623516+07	25
75	451425	t	50	Stok opname bulanan	2026-05-12 20:23:22.629696+07	25
76	238299	t	34	Stok opname bulanan	2026-05-31 20:23:22.636232+07	26
77	236966	f	15	Stok keluar penjualan	2026-05-17 20:23:22.64225+07	26
78	233707	t	85	Stok keluar penjualan	2026-06-01 20:23:22.648596+07	26
79	411445	t	59	Retur barang	2026-05-31 20:23:22.654628+07	27
80	411557	f	47	Stok keluar penjualan	2026-06-08 20:23:22.66016+07	27
81	418961	t	33	Stok keluar penjualan	2026-06-05 20:23:22.666906+07	27
82	98515	t	88	Stok opname bulanan	2026-05-13 20:23:22.672706+07	28
83	105967	f	94	Retur barang	2026-06-06 20:23:22.678262+07	28
84	105664	t	60	Stok masuk dari supplier	2026-06-05 20:23:22.684055+07	28
85	237223	t	85	Stok opname bulanan	2026-06-07 20:23:22.689688+07	29
86	238928	f	87	Retur barang	2026-05-11 20:23:22.695037+07	29
87	237832	t	54	Stok masuk dari supplier	2026-05-18 20:23:22.701339+07	29
88	166323	t	91	Stok keluar penjualan	2026-06-07 20:23:22.707532+07	30
89	171631	f	91	Stok opname bulanan	2026-06-08 20:23:22.713336+07	30
90	165208	t	41	Retur barang	2026-05-19 20:23:22.718776+07	30
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."user" (id_user, public_id, username, email, password, nama_lengkap, foto_profil, id_role) FROM stdin;
6	d45067d3-6272-4fcb-8cdd-b5b694fbbb19	customer694	customer444@example.com	$2a$12$RTBOTiHQZQ8regxh8oJkOOFp2rpHBc/GqnF1teFedYKSSv0tse./W	Customer Baru		1
7	8c7b9f6b-227b-4103-a2e8-47752ab5a776	hamim_customer	customer@mantra.com	$2a$10$OpdCPFJhDib2ELb6llSA1.iadWeHXU0pMs9eCVeDmXxJc60D.UChC	Rajaba Hamim	https://ui-avatars.com/api/?name=Rajaba+Hamim&background=3b82f6&color=fff	1
8	4ccbf168-3081-4271-8675-4a3dead4c6e4	karyawan199	karyawan30@mantra.web.id	$2a$12$WvkMX8jBa/CWEdiDadaGkOax8H5MHpeugHBmwmpp0FPCT7tWfsiW6	Karyawan Baru		3
3	dc37ad63-9dee-4144-9bf0-8a3e9359c336	riztika_kurir	kurir@mantra.com	$2a$10$aB/Ri2ZhhFITeKG0nTssheHys6lOJxeonzjrCMMzePw2Io7jqrt0y	Karyawan Updated	https://ui-avatars.com/api/?name=Riztika+Merizta&background=f59e0b&color=fff	4
1	a0a1416a-0528-4f50-adfc-dba4b77c0210	terra_admin	terra_admin@mantra.web.id	$2a$10$aEHQbsyGADvLu4xGn6cAeesFxHM8LF81IoAEiQvu3unNsRYsuWhfC	Admin Updated	https://ui-avatars.com/api/?name=Terra+Surya&background=6366f1&color=fff	2
9	691c2c6d-3b8f-4827-9b56-cfcb2ebb5a25	customer809	customer248@example.com	$2a$12$UOLqgi1c0rfM31ukMU7hjuR2.1NmSW5W7vJ0kfCDLsbannQWQ/Kqy	Customer Baru		1
4	06e57a0a-06d8-42d2-b896-ae51826aa7c1	hamim_customer	admin@mantra.com	$2a$10$er9AXZum2xqXrel6WCA/JeZ6M9P9dwxDJ6y.FvKbbAqV3aP/36Kqy	John Updated	https://ui-avatars.com/api/?name=Rajaba+Hamim&background=3b82f6&color=fff	1
2	c0371782-5673-42e8-9b7a-62730fa8705d	nabila_kasir	kasir@mantra.com	$2a$10$toflq/RIp3VGi3unEU5Px.w7DOHEJiLMYXxs7cQ6cCP5rQbuX7OWS	Nabila Az Zahra	https://storage.mantra.web.id/mantra-storage/karyawan/2026/06/23bfe060-45e2-4f8b-aa8c-fa46dffb3377.png	3
5	124743c8-b610-4190-b9f7-32daea75c79a	karyawan608	karyawan496@mantra.web.id	$2a$12$AyW8jSl6ocb40FhvNjpi1.bfSffBnvtH/c0MKKzMoDPfxtJNpXt8W	Karyawan Baru	https://storage.mantra.web.id/mantra-storage/karyawan/2026/06/3eeb8538-d8bb-45be-a409-2667ea9f9d1e.png	3
10	9dffe424-e904-42a2-9492-536aa31996b4	customer356	customer122@example.com	$2a$12$Kp67uTkQehhwedLHrF7g1.2dpy23Z6rsvDTzxIa855R5gpvVBxUCa	Customer Baru		1
\.


--
-- Name: alamat_id_alamat_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.alamat_id_alamat_seq', 6, true);


--
-- Name: barang_id_barang_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.barang_id_barang_seq', 17, true);


--
-- Name: barcode_id_barcode_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.barcode_id_barcode_seq', 324, true);


--
-- Name: customer_id_customer_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.customer_id_customer_seq', 5, true);


--
-- Name: detail_pembayaran_id_detail_pembayaran_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.detail_pembayaran_id_detail_pembayaran_seq', 1, false);


--
-- Name: detail_pesanan_id_detail_pesanan_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.detail_pesanan_id_detail_pesanan_seq', 646, true);


--
-- Name: detail_spesifikasi_id_detail_spesifikasi_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.detail_spesifikasi_id_detail_spesifikasi_seq', 14, true);


--
-- Name: diskon_id_diskon_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.diskon_id_diskon_seq', 6, true);


--
-- Name: ekspedisi_id_ekspedisi_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ekspedisi_id_ekspedisi_seq', 8, true);


--
-- Name: ekspedisi_layanan_id_ekspedisi_layanan_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ekspedisi_layanan_id_ekspedisi_layanan_seq', 13, true);


--
-- Name: karyawan_id_karyawan_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.karyawan_id_karyawan_seq', 4, true);


--
-- Name: kasir_id_kasir_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kasir_id_kasir_seq', 3, true);


--
-- Name: kategori_id_kategori_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kategori_id_kategori_seq', 11, true);


--
-- Name: keranjang_id_keranjang_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.keranjang_id_keranjang_seq', 11, true);


--
-- Name: kurir_id_kurir_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.kurir_id_kurir_seq', 1, true);


--
-- Name: metode_pembayaran_id_metode_pembayaran_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.metode_pembayaran_id_metode_pembayaran_seq', 7, true);


--
-- Name: notifikasi_id_notifikasi_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifikasi_id_notifikasi_seq', 10, true);


--
-- Name: pembayaran_id_pembayaran_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pembayaran_id_pembayaran_seq', 233, true);


--
-- Name: pengantaran_id_pengantaran_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pengantaran_id_pengantaran_seq', 142, true);


--
-- Name: pesanan_id_pesanan_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.pesanan_id_pesanan_seq', 230, true);


--
-- Name: refresh_token_id_token_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.refresh_token_id_token_seq', 26, true);


--
-- Name: role_id_role_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.role_id_role_seq', 4, true);


--
-- Name: satuan_id_satuan_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.satuan_id_satuan_seq', 4, true);


--
-- Name: spesifikasi_barang_id_spesifikasi_barang_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.spesifikasi_barang_id_spesifikasi_barang_seq', 32, true);


--
-- Name: spesifikasi_id_spesifikasi_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.spesifikasi_id_spesifikasi_seq', 2, true);


--
-- Name: status_pengantaran_id_status_pengantaran_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.status_pengantaran_id_status_pengantaran_seq', 5, true);


--
-- Name: stok_opname_id_stok_opname_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stok_opname_id_stok_opname_seq', 94, true);


--
-- Name: user_id_user_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id_user_seq', 10, true);


--
-- Name: alamat alamat_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.alamat
    ADD CONSTRAINT alamat_pkey PRIMARY KEY (id_alamat);


--
-- Name: barang barang_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barang
    ADD CONSTRAINT barang_pkey PRIMARY KEY (id_barang);


--
-- Name: barcode barcode_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barcode
    ADD CONSTRAINT barcode_pkey PRIMARY KEY (id_barcode);


--
-- Name: customer customer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT customer_pkey PRIMARY KEY (id_customer);


--
-- Name: detail_pembayaran detail_pembayaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pembayaran
    ADD CONSTRAINT detail_pembayaran_pkey PRIMARY KEY (id_detail_pembayaran);


--
-- Name: detail_pesanan detail_pesanan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pesanan
    ADD CONSTRAINT detail_pesanan_pkey PRIMARY KEY (id_detail_pesanan);


--
-- Name: detail_spesifikasi detail_spesifikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_spesifikasi
    ADD CONSTRAINT detail_spesifikasi_pkey PRIMARY KEY (id_detail_spesifikasi);


--
-- Name: diskon diskon_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.diskon
    ADD CONSTRAINT diskon_pkey PRIMARY KEY (id_diskon);


--
-- Name: ekspedisi_layanan ekspedisi_layanan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ekspedisi_layanan
    ADD CONSTRAINT ekspedisi_layanan_pkey PRIMARY KEY (id_ekspedisi_layanan);


--
-- Name: ekspedisi ekspedisi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ekspedisi
    ADD CONSTRAINT ekspedisi_pkey PRIMARY KEY (id_ekspedisi);


--
-- Name: karyawan karyawan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan
    ADD CONSTRAINT karyawan_pkey PRIMARY KEY (id_karyawan);


--
-- Name: kasir kasir_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT kasir_pkey PRIMARY KEY (id_kasir);


--
-- Name: kategori kategori_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kategori
    ADD CONSTRAINT kategori_pkey PRIMARY KEY (id_kategori);


--
-- Name: keranjang keranjang_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keranjang
    ADD CONSTRAINT keranjang_pkey PRIMARY KEY (id_keranjang);


--
-- Name: kurir kurir_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT kurir_pkey PRIMARY KEY (id_kurir);


--
-- Name: metode_pembayaran metode_pembayaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.metode_pembayaran
    ADD CONSTRAINT metode_pembayaran_pkey PRIMARY KEY (id_metode_pembayaran);


--
-- Name: notifikasi notifikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT notifikasi_pkey PRIMARY KEY (id_notifikasi);


--
-- Name: pembayaran pembayaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pembayaran
    ADD CONSTRAINT pembayaran_pkey PRIMARY KEY (id_pembayaran);


--
-- Name: pengantaran pengantaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran
    ADD CONSTRAINT pengantaran_pkey PRIMARY KEY (id_pengantaran);


--
-- Name: pesanan pesanan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT pesanan_pkey PRIMARY KEY (id_pesanan);


--
-- Name: refresh_token refresh_token_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT refresh_token_pkey PRIMARY KEY (id_token);


--
-- Name: role role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT role_pkey PRIMARY KEY (id_role);


--
-- Name: satuan satuan_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.satuan
    ADD CONSTRAINT satuan_pkey PRIMARY KEY (id_satuan);


--
-- Name: spesifikasi_barang spesifikasi_barang_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi_barang
    ADD CONSTRAINT spesifikasi_barang_pkey PRIMARY KEY (id_spesifikasi_barang);


--
-- Name: spesifikasi spesifikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi
    ADD CONSTRAINT spesifikasi_pkey PRIMARY KEY (id_spesifikasi);


--
-- Name: status_pengantaran status_pengantaran_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.status_pengantaran
    ADD CONSTRAINT status_pengantaran_pkey PRIMARY KEY (id_status_pengantaran);


--
-- Name: stok_opname stok_opname_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stok_opname
    ADD CONSTRAINT stok_opname_pkey PRIMARY KEY (id_stok_opname);


--
-- Name: customer uni_customer_id_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT uni_customer_id_user UNIQUE (id_user);


--
-- Name: karyawan uni_karyawan_id_user; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan
    ADD CONSTRAINT uni_karyawan_id_user UNIQUE (id_user);


--
-- Name: karyawan uni_karyawan_nik; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan
    ADD CONSTRAINT uni_karyawan_nik UNIQUE (nik);


--
-- Name: karyawan uni_karyawan_no_telp; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan
    ADD CONSTRAINT uni_karyawan_no_telp UNIQUE (no_telp);


--
-- Name: kasir uni_kasir_id_karyawan; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT uni_kasir_id_karyawan UNIQUE (id_karyawan);


--
-- Name: kurir uni_kurir_id_karyawan; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT uni_kurir_id_karyawan UNIQUE (id_karyawan);


--
-- Name: role uni_role_nama_role; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role
    ADD CONSTRAINT uni_role_nama_role UNIQUE (nama_role);


--
-- Name: user uni_user_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT uni_user_email UNIQUE (email);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id_user);


--
-- Name: idx_alamat_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_alamat_public_id ON public.alamat USING btree (public_id);


--
-- Name: idx_barang_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_barang_public_id ON public.barang USING btree (public_id);


--
-- Name: idx_customer_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_customer_public_id ON public.customer USING btree (public_id);


--
-- Name: idx_detail_pembayaran_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_detail_pembayaran_public_id ON public.detail_pembayaran USING btree (public_id);


--
-- Name: idx_diskon_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_diskon_public_id ON public.diskon USING btree (public_id);


--
-- Name: idx_ekspedisi_layanan_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_ekspedisi_layanan_public_id ON public.ekspedisi_layanan USING btree (public_id);


--
-- Name: idx_ekspedisi_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_ekspedisi_public_id ON public.ekspedisi USING btree (public_id);


--
-- Name: idx_karyawan_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_karyawan_public_id ON public.karyawan USING btree (public_id);


--
-- Name: idx_kasir_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_kasir_public_id ON public.kasir USING btree (public_id);


--
-- Name: idx_kategori_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_kategori_public_id ON public.kategori USING btree (public_id);


--
-- Name: idx_keranjang_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_keranjang_public_id ON public.keranjang USING btree (public_id);


--
-- Name: idx_kurir_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_kurir_public_id ON public.kurir USING btree (public_id);


--
-- Name: idx_metode_pembayaran_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_metode_pembayaran_public_id ON public.metode_pembayaran USING btree (public_id);


--
-- Name: idx_pengantaran_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_pengantaran_public_id ON public.pengantaran USING btree (public_id);


--
-- Name: idx_pesanan_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_pesanan_public_id ON public.pesanan USING btree (public_id);


--
-- Name: idx_user_public_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_user_public_id ON public."user" USING btree (public_id);


--
-- Name: alamat fk_alamat_customer; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.alamat
    ADD CONSTRAINT fk_alamat_customer FOREIGN KEY (id_customer) REFERENCES public.customer(id_customer);


--
-- Name: barang fk_barang_diskon; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barang
    ADD CONSTRAINT fk_barang_diskon FOREIGN KEY (id_diskon) REFERENCES public.diskon(id_diskon);


--
-- Name: barang fk_barang_kategori; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barang
    ADD CONSTRAINT fk_barang_kategori FOREIGN KEY (id_kategori) REFERENCES public.kategori(id_kategori);


--
-- Name: barang fk_barang_satuan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barang
    ADD CONSTRAINT fk_barang_satuan FOREIGN KEY (id_satuan) REFERENCES public.satuan(id_satuan);


--
-- Name: barcode fk_barcode_spesifikasi_barang; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.barcode
    ADD CONSTRAINT fk_barcode_spesifikasi_barang FOREIGN KEY (id_spesifikasi_barang) REFERENCES public.spesifikasi_barang(id_spesifikasi_barang);


--
-- Name: customer fk_customer_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer
    ADD CONSTRAINT fk_customer_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- Name: detail_pesanan fk_detail_pesanan_spesifikasi_barang; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pesanan
    ADD CONSTRAINT fk_detail_pesanan_spesifikasi_barang FOREIGN KEY (id_spesifikasi_barang) REFERENCES public.spesifikasi_barang(id_spesifikasi_barang);


--
-- Name: detail_spesifikasi fk_detail_spesifikasi_spesifikasi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_spesifikasi
    ADD CONSTRAINT fk_detail_spesifikasi_spesifikasi FOREIGN KEY (id_spesifikasi) REFERENCES public.spesifikasi(id_spesifikasi);


--
-- Name: ekspedisi_layanan fk_ekspedisi_layanan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ekspedisi_layanan
    ADD CONSTRAINT fk_ekspedisi_layanan FOREIGN KEY (id_ekspedisi) REFERENCES public.ekspedisi(id_ekspedisi) ON DELETE CASCADE;


--
-- Name: karyawan fk_karyawan_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.karyawan
    ADD CONSTRAINT fk_karyawan_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- Name: kasir fk_kasir_karyawan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kasir
    ADD CONSTRAINT fk_kasir_karyawan FOREIGN KEY (id_karyawan) REFERENCES public.karyawan(id_karyawan);


--
-- Name: keranjang fk_keranjang_customer; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keranjang
    ADD CONSTRAINT fk_keranjang_customer FOREIGN KEY (id_customer) REFERENCES public.customer(id_customer);


--
-- Name: keranjang fk_keranjang_spesifikasi_barang; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keranjang
    ADD CONSTRAINT fk_keranjang_spesifikasi_barang FOREIGN KEY (id_spesifikasi_barang) REFERENCES public.spesifikasi_barang(id_spesifikasi_barang);


--
-- Name: kurir fk_kurir_karyawan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kurir
    ADD CONSTRAINT fk_kurir_karyawan FOREIGN KEY (id_karyawan) REFERENCES public.karyawan(id_karyawan);


--
-- Name: notifikasi fk_notifikasi_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifikasi
    ADD CONSTRAINT fk_notifikasi_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- Name: detail_pembayaran fk_pembayaran_detail_pembayaran; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pembayaran
    ADD CONSTRAINT fk_pembayaran_detail_pembayaran FOREIGN KEY (id_pembayaran) REFERENCES public.pembayaran(id_pembayaran);


--
-- Name: pembayaran fk_pembayaran_metode_pembayaran; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pembayaran
    ADD CONSTRAINT fk_pembayaran_metode_pembayaran FOREIGN KEY (id_metode_pembayaran) REFERENCES public.metode_pembayaran(id_metode_pembayaran);


--
-- Name: pembayaran fk_pembayaran_pesanan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pembayaran
    ADD CONSTRAINT fk_pembayaran_pesanan FOREIGN KEY (id_pesanan) REFERENCES public.pesanan(id_pesanan);


--
-- Name: pengantaran fk_pengantaran_ekspedisi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran
    ADD CONSTRAINT fk_pengantaran_ekspedisi FOREIGN KEY (id_ekspedisi) REFERENCES public.ekspedisi(id_ekspedisi);


--
-- Name: pengantaran fk_pengantaran_kurir; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran
    ADD CONSTRAINT fk_pengantaran_kurir FOREIGN KEY (id_kurir) REFERENCES public.kurir(id_kurir);


--
-- Name: pengantaran fk_pengantaran_pesanan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran
    ADD CONSTRAINT fk_pengantaran_pesanan FOREIGN KEY (id_pesanan) REFERENCES public.pesanan(id_pesanan);


--
-- Name: pengantaran fk_pengantaran_status_pengantaran; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pengantaran
    ADD CONSTRAINT fk_pengantaran_status_pengantaran FOREIGN KEY (id_status_pengantaran) REFERENCES public.status_pengantaran(id_status_pengantaran);


--
-- Name: pesanan fk_pesanan_alamat; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT fk_pesanan_alamat FOREIGN KEY (id_alamat) REFERENCES public.alamat(id_alamat);


--
-- Name: pesanan fk_pesanan_customer; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT fk_pesanan_customer FOREIGN KEY (id_customer) REFERENCES public.customer(id_customer);


--
-- Name: detail_pesanan fk_pesanan_detail_pesanan; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detail_pesanan
    ADD CONSTRAINT fk_pesanan_detail_pesanan FOREIGN KEY (id_pesanan) REFERENCES public.pesanan(id_pesanan);


--
-- Name: pesanan fk_pesanan_ekspedisi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT fk_pesanan_ekspedisi FOREIGN KEY (id_ekspedisi) REFERENCES public.ekspedisi(id_ekspedisi);


--
-- Name: pesanan fk_pesanan_kasir; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT fk_pesanan_kasir FOREIGN KEY (id_kasir) REFERENCES public.kasir(id_kasir);


--
-- Name: pesanan fk_pesanan_layanan_ekspedisi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pesanan
    ADD CONSTRAINT fk_pesanan_layanan_ekspedisi FOREIGN KEY (id_layanan_ekspedisi) REFERENCES public.ekspedisi_layanan(id_ekspedisi_layanan);


--
-- Name: refresh_token fk_refresh_token_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT fk_refresh_token_user FOREIGN KEY (id_user) REFERENCES public."user"(id_user);


--
-- Name: spesifikasi_barang fk_spesifikasi_barang_barang; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi_barang
    ADD CONSTRAINT fk_spesifikasi_barang_barang FOREIGN KEY (id_barang) REFERENCES public.barang(id_barang);


--
-- Name: spesifikasi_barang fk_spesifikasi_barang_detail_spesifikasi; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spesifikasi_barang
    ADD CONSTRAINT fk_spesifikasi_barang_detail_spesifikasi FOREIGN KEY (id_detail_spesifikasi) REFERENCES public.detail_spesifikasi(id_detail_spesifikasi);


--
-- Name: stok_opname fk_stok_opname_spesifikasi_barang; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stok_opname
    ADD CONSTRAINT fk_stok_opname_spesifikasi_barang FOREIGN KEY (id_spesifikasi_barang) REFERENCES public.spesifikasi_barang(id_spesifikasi_barang);


--
-- Name: user fk_user_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT fk_user_role FOREIGN KEY (id_role) REFERENCES public.role(id_role);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


--
-- PostgreSQL database dump complete
--

\unrestrict lFqV4I3g6gMjLcWOkldox15LDX5VkhhvDl9dIAgfTh6wQKyn4q7nBIj7kKyNC3h

