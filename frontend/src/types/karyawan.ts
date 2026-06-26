export interface Karyawan {
  id: number;
  kode: string;
  nama: string;
  tanggal_mulai: string;
  tanggal_habis: string;
  created_at: string;
  updated_at: string;
}

export interface KaryawanRequest {
  kode: string;
  nama: string;
  tanggal_mulai: string;
  tanggal_habis: string;
}

export interface ApiResponse<T = unknown> {
  message: string;
  data?: T;
  errors?: Record<string, string> | string;
}
