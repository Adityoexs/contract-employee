import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { karyawanService } from '../services/karyawanService';
import type { Karyawan } from '../types/karyawan';

function formatDate(iso: string): string {
  if (!iso) return '-';
  return iso.slice(0, 10);
}

export default function DetailPage() {
  const { id } = useParams<{ id: string }>();
  const [data, setData] = useState<Karyawan | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    karyawanService
      .getById(Number(id))
      .then((res) => setData(res.data ?? null))
      .catch(() => toast.error('Data tidak ditemukan'))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="container"><p className="loading">Memuat data...</p></div>;
  if (!data) return <div className="container"><p className="empty">Data tidak ditemukan.</p></div>;

  return (
    <div className="container">
      <div className="page-header">
        <h1>Detail Karyawan Kontrak</h1>
        <Link to="/" className="btn btn-secondary">
          &larr; Kembali
        </Link>
      </div>

      <div className="detail-card">
        <div className="detail-row">
          <span className="detail-label">Kode</span>
          <span className="detail-value">{data.kode}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Nama</span>
          <span className="detail-value">{data.nama}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Tanggal Mulai</span>
          <span className="detail-value">{formatDate(data.tanggal_mulai)}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Tanggal Habis</span>
          <span className="detail-value">{formatDate(data.tanggal_habis)}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Dibuat</span>
          <span className="detail-value">{new Date(data.created_at).toLocaleString('id-ID')}</span>
        </div>
        <div className="detail-row">
          <span className="detail-label">Diperbarui</span>
          <span className="detail-value">{new Date(data.updated_at).toLocaleString('id-ID')}</span>
        </div>
      </div>

      <div style={{ marginTop: '1.5rem' }}>
        <Link to={`/edit/${data.id}`} className="btn btn-warning">
          Edit Data
        </Link>
      </div>
    </div>
  );
}
