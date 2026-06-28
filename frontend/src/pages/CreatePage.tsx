import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { isAxiosError } from 'axios';
import KaryawanForm from '../components/KaryawanForm';
import { karyawanService } from '../services/karyawanService';
import type { KaryawanCreateRequest } from '../types/karyawan';

export default function CreatePage() {
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [serverErrors, setServerErrors] = useState<Record<string, string> | undefined>();

  const handleSubmit = async (data: KaryawanCreateRequest) => {
    setIsSubmitting(true);
    setServerErrors(undefined);
    try {
      await karyawanService.create(data);
      toast.success('Data karyawan berhasil ditambahkan');
      navigate('/');
    } catch (err) {
      if (isAxiosError(err) && err.response) {
        const body = err.response.data;
        if (typeof body.errors === 'object' && body.errors !== null) {
          setServerErrors(body.errors as Record<string, string>);
          toast.error(body.message ?? 'Validasi gagal');
        } else {
          toast.error(body.message ?? 'Gagal menyimpan data');
        }
      } else {
        toast.error('Terjadi kesalahan, coba lagi');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="container">
      <div className="page-header">
        <h1>Tambah Karyawan Kontrak</h1>
        <Link to="/" className="btn btn-secondary">
          &larr; Kembali
        </Link>
      </div>

      <div className="form-card">
        <KaryawanForm
          onSubmit={handleSubmit}
          serverErrors={serverErrors}
          isSubmitting={isSubmitting}
          submitLabel="Tambah Karyawan"
        />
      </div>
    </div>
  );
}
