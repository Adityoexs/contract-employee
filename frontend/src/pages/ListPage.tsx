import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { karyawanService } from '../services/karyawanService';
import type { Karyawan } from '../types/karyawan';

function formatDate(iso: string): string {
  if (!iso) return '-';
  return iso.slice(0, 10);
}

export default function ListPage() {
  const [list, setList] = useState<Karyawan[]>([]);
  const [loading, setLoading] = useState(true);
  const [importing, setImporting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await karyawanService.getAll();
      setList(res.data ?? []);
    } catch {
      toast.error('Gagal memuat data karyawan');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleDelete = async (id: number, nama: string) => {
    if (!window.confirm(`Hapus karyawan "${nama}"?`)) return;
    try {
      await karyawanService.delete(id);
      toast.success('Data berhasil dihapus');
      fetchData();
    } catch {
      toast.error('Gagal menghapus data');
    }
  };

  const handleDownloadTemplate = async () => {
    try {
      await karyawanService.downloadTemplate();
      toast.success('Template berhasil diunduh');
    } catch {
      toast.error('Gagal mengunduh template');
    }
  };

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Reset input so the same file can be re-selected if needed.
    e.target.value = '';

    setImporting(true);
    try {
      const res = await karyawanService.importExcel(file);
      const inserted = res.data?.inserted ?? 0;
      toast.success(`Import berhasil: ${inserted} data ditambahkan`);
      fetchData();
    } catch (err: unknown) {
      // Try to surface validation errors from the backend response.
      const axiosErr = err as { response?: { data?: { message?: string; errors?: unknown } } };
      const detail = axiosErr?.response?.data;
      if (detail?.errors && Array.isArray(detail.errors)) {
        const msgs = (detail.errors as Array<{ row: number; field: string; message: string }>)
          .slice(0, 3)
          .map((e) => `Baris ${e.row} (${e.field}): ${e.message}`)
          .join('\n');
        toast.error(`Import gagal:\n${msgs}`, { duration: 6000 });
      } else {
        toast.error(detail?.message ?? 'Import gagal');
      }
    } finally {
      setImporting(false);
    }
  };

  return (
    <div className="container">
      <div className="page-header">
        <h1>Daftar Karyawan Kontrak</h1>
        <div className="actions">
          <button
            type="button"
            className="btn btn-secondary"
            onClick={handleDownloadTemplate}
          >
            ⬇ Template Excel
          </button>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={handleImportClick}
            disabled={importing}
          >
            {importing ? 'Mengimpor...' : '⬆ Import Excel'}
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept=".xlsx,.xls"
            style={{ display: 'none' }}
            onChange={handleFileChange}
          />
          <Link to="/create" className="btn btn-primary">
            + Tambah Karyawan
          </Link>
        </div>
      </div>

      {loading ? (
        <p className="loading">Memuat data...</p>
      ) : list.length === 0 ? (
        <p className="empty">Belum ada data karyawan kontrak.</p>
      ) : (
        <div className="table-wrapper">
          <table className="table">
            <thead>
              <tr>
                <th>No</th>
                <th>Kode</th>
                <th>Nama</th>
                <th>Tanggal Mulai</th>
                <th>Tanggal Habis</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {list.map((k, idx) => (
                <tr key={k.id}>
                  <td>{idx + 1}</td>
                  <td>{k.kode}</td>
                  <td>{k.nama}</td>
                  <td>{formatDate(k.tanggal_mulai)}</td>
                  <td>{formatDate(k.tanggal_habis)}</td>
                  <td className="actions">
                    <Link to={`/detail/${k.id}`} className="btn btn-sm btn-info">
                      Detail
                    </Link>
                    <Link to={`/edit/${k.id}`} className="btn btn-sm btn-warning">
                      Edit
                    </Link>
                    <button
                      className="btn btn-sm btn-danger"
                      onClick={() => handleDelete(k.id, k.nama)}
                    >
                      Hapus
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
