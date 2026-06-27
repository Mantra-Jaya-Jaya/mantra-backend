package utils

import (
	"backend-mantra/config"
	"log"
	"time"
)

// StartOrderAutoCompletionCron runs in the background to automatically complete orders
// that have been "Dikirim" for more than 24 hours.
func StartOrderAutoCompletionCron() {
	go func() {
		for {
			// Get status IDs
			dikirimID := GetStatusPesananIDSafe("Dikirim")
			selesaiID := GetStatusPesananIDSafe("Selesai")
			internalKurirID := GetTipeKurirIDSafe("internal")
			tibaID := GetStatusPengantaranIDSafe("Tiba di Tujuan")

			if dikirimID != 0 && selesaiID != 0 && internalKurirID != 0 && tibaID != 0 {
				threshold := time.Now().Add(-24 * time.Hour)

				// Update all matching orders
				// 1. Tipe kurir internal
				// 2. Status pesanan = Dikirim
				// 3. Pengantaran has status = Tiba di Tujuan
				// 4. Pengantaran.FotoBuktiPengiriman is NOT NULL and NOT ''
				// 5. Pengantaran.WaktuSampai < threshold

				result := config.DB.Exec(`
					UPDATE pesanan 
					SET id_status_pesanan = ? 
					FROM pengantaran 
					WHERE pesanan.id_pesanan = pengantaran.id_pesanan 
					AND pesanan.id_status_pesanan = ? 
					AND pesanan.id_tipe_kurir = ? 
					AND pengantaran.id_status_pengantaran = ? 
					AND pengantaran.foto_bukti_pengiriman IS NOT NULL 
					AND pengantaran.foto_bukti_pengiriman != '' 
					AND pengantaran.waktu_sampai < ?`,
					selesaiID, dikirimID, internalKurirID, tibaID, threshold)

				if result.Error != nil {
					log.Printf("[CRON] Error auto-completing orders: %v", result.Error)
				} else if result.RowsAffected > 0 {
					log.Printf("[CRON] Auto-completed %d orders.", result.RowsAffected)
				}
			}

			// Run every 1 hour
			time.Sleep(1 * time.Hour)
		}
	}()
}
