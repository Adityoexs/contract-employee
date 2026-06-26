import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { isAxiosError } from 'axios';
import KaryawanForm from '../components/KaryawanForm';
import { karyawanService } from '../services/karyawanService';
import type { Karyawan, KaryawanRequest } from '../types/karyawan';

export default function EditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [data, setData] = useState<Karyawan | null>(null);
  const [loading, setLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [serverErrors, setServerErrors] = useState<Record<string, string> | undefined>();

  useEffect(() => {
    if (!id) return;
    karyawanService
      .getById(Number(id))
      .then((res) => setData(res.data ?? null))
      .catch(() => toast.error('Data tidak ditemukan'))
      .finally(() => setLoading(false));
  }, [id]);

  const handleSubmit = async (formData: KaryawanRequest) => {
    if (!id) return;
    setIsSubmitting(true);
    setServerErrors(undefined);
    try {
      await karyawanService.update(Number(id), formData);
      toast.success('Data karyawan berhasil diperbarui');
      navigate('/');
    } catch (err) {
      if (isAxiosError(err) && err.response) {
        const body = err.response.data;
        if (typeof body.errors === 'object' && body.errors !== null) {
          setServerErrors(body.errors as Record<string, string>);
          toast.error(body.message ?? 'Validasi gagal');
        } else {
          toast.error(body.message ?? 'Gagal memperbarui data');
        }
      } else {
        toast.error('Terjadi kesalahan, coba lagi');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading)
    return (
      <div className="container">
        <p className="loading">Memuat data...</p>
      </div>
    );

  if (!data)
    return (
      <div className="container">
        <p className="empty">Data tidak ditemukan.</p>
        <Link to="/" className="btn btn-secondary">
          Kembali
        </Link>
      </div>
    );

  return (
    <div className="container">
      <div className="page-header">
        <h1>Edit Karyawan Kontrak</h1>
        <Link to="/" className="btn btn-secondary">
          &larr; Kembali
        </Link>
      </div>

      <div className="form-card">
        <KaryawanForm
          defaultValues={{
            kode: data.kode,
            nama: data.nama,
            tanggal_mulai: data.tanggal_mulai.slice(0, 10),
            tanggal_habis: data.tanggal_habis.slice(0, 10),
          }}
          onSubmit={handleSubmit}
          serverErrors={serverErrors}
          isSubmitting={isSubmitting}
          submitLabel="Perbarui Karyawan"
        />
      </div>
    </div>
  );
}
