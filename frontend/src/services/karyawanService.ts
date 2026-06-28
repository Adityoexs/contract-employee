import axios from 'axios';
import type { ApiResponse, Karyawan, KaryawanRequest } from '../types/karyawan';

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api',
  headers: { 'Content-Type': 'application/json' },
});

export const karyawanService = {
  getAll(): Promise<ApiResponse<Karyawan[]>> {
    return http.get<ApiResponse<Karyawan[]>>('/karyawan').then((r) => r.data);
  },

  getById(id: number): Promise<ApiResponse<Karyawan>> {
    return http.get<ApiResponse<Karyawan>>(`/karyawan/${id}`).then((r) => r.data);
  },

  create(data: KaryawanRequest): Promise<ApiResponse<Karyawan>> {
    return http.post<ApiResponse<Karyawan>>('/karyawan', data).then((r) => r.data);
  },

  update(id: number, data: KaryawanRequest): Promise<ApiResponse<Karyawan>> {
    return http
      .put<ApiResponse<Karyawan>>(`/karyawan/${id}`, data)
      .then((r) => r.data);
  },

  delete(id: number): Promise<ApiResponse> {
    return http.delete<ApiResponse>(`/karyawan/${id}`).then((r) => r.data);
  },

  downloadTemplate(): Promise<void> {
    return http
      .get('/karyawan/import/template', { responseType: 'blob' })
      .then((r) => {
        const url = URL.createObjectURL(r.data as Blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'template_karyawan.xlsx';
        a.click();
        URL.revokeObjectURL(url);
      });
  },

  importExcel(file: File): Promise<ApiResponse<{ inserted: number }>> {
    const form = new FormData();
    form.append('file', file);
    return http
      .post<ApiResponse<{ inserted: number }>>('/karyawan/import', form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then((r) => r.data);
  },
};
