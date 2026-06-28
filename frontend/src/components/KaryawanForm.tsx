import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import type { KaryawanCreateRequest } from '../types/karyawan';

const schema = z
  .object({
    nama: z
      .string()
      .min(1, 'Nama wajib diisi')
      .refine((v) => v.trim().length > 0, 'Nama tidak boleh hanya spasi'),
    tanggal_mulai: z
      .string()
      .min(1, 'Tanggal mulai wajib diisi')
      .regex(/^\d{4}-\d{2}-\d{2}$/, 'Format tanggal tidak valid (YYYY-MM-DD)'),
    tanggal_habis: z
      .string()
      .min(1, 'Tanggal habis wajib diisi')
      .regex(/^\d{4}-\d{2}-\d{2}$/, 'Format tanggal tidak valid (YYYY-MM-DD)'),
  })
  .refine((d) => new Date(d.tanggal_habis) >= new Date(d.tanggal_mulai), {
    message: 'Tanggal habis harus lebih besar atau sama dengan tanggal mulai',
    path: ['tanggal_habis'],
  });

type FormValues = z.infer<typeof schema>;

interface Props {
  defaultValues?: Partial<KaryawanCreateRequest>;
  onSubmit: (data: KaryawanCreateRequest) => Promise<void>;
  serverErrors?: Record<string, string>;
  isSubmitting?: boolean;
  submitLabel?: string;
}

export default function KaryawanForm({
  defaultValues,
  onSubmit,
  serverErrors,
  isSubmitting = false,
  submitLabel = 'Simpan',
}: Props) {
  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      nama: defaultValues?.nama ?? '',
      tanggal_mulai: defaultValues?.tanggal_mulai?.slice(0, 10) ?? '',
      tanggal_habis: defaultValues?.tanggal_habis?.slice(0, 10) ?? '',
    },
  });

  // Sync server-side validation errors into react-hook-form
  useEffect(() => {
    if (!serverErrors) return;
    (Object.entries(serverErrors) as [keyof FormValues, string][]).forEach(
      ([field, message]) => {
        setError(field, { type: 'server', message });
      }
    );
  }, [serverErrors, setError]);

  return (
    <form onSubmit={handleSubmit(onSubmit)} noValidate className="form">
      <div className="form-group">
        <label htmlFor="nama">Nama <span className="required">*</span></label>
        <input id="nama" type="text" placeholder="Nama lengkap karyawan" {...register('nama')} />
        {errors.nama && <p className="field-error">{errors.nama.message}</p>}
      </div>

      <div className="form-group">
        <label htmlFor="tanggal_mulai">Tanggal Mulai <span className="required">*</span></label>
        <input id="tanggal_mulai" type="date" {...register('tanggal_mulai')} />
        {errors.tanggal_mulai && (
          <p className="field-error">{errors.tanggal_mulai.message}</p>
        )}
      </div>

      <div className="form-group">
        <label htmlFor="tanggal_habis">Tanggal Habis <span className="required">*</span></label>
        <input id="tanggal_habis" type="date" {...register('tanggal_habis')} />
        {errors.tanggal_habis && (
          <p className="field-error">{errors.tanggal_habis.message}</p>
        )}
      </div>

      <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
        {isSubmitting ? 'Menyimpan...' : submitLabel}
      </button>
    </form>
  );
}
