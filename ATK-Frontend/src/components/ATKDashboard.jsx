import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Plus, Pencil, Trash2, Package, X } from 'lucide-react';
import Swal from 'sweetalert2';

const ATKDashboard = () => {
  const [atkItems, setAtkItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [currentItem, setCurrentItem] = useState({ id: 0, nama: '', jenis: '', qty: 0 });
  const [error, setError] = useState('');

  // Fetch all ATK items
  const fetchATKItems = async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/atk');
      setAtkItems(response.data || []);
      setError('');
    } catch (err) {
      console.error('Error fetching ATK items:', err);
      setError('Gagal mengambil data ATK. Pastikan backend berjalan di http://localhost:5200');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchATKItems();
  }, []);

  // Open modal for create
  const handleCreate = () => {
    setEditMode(false);
    setCurrentItem({ id: 0, nama: '', jenis: '', qty: 0 });
    setShowModal(true);
  };

  // Open modal for edit
  const handleEdit = (item) => {
    setEditMode(true);
    setCurrentItem(item);
    setShowModal(true);
  };

  // Handle form submit (create or update)
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    try {
      if (editMode) {
        // Update existing item
        await axios.put(`/api/atk/${currentItem.id}`, {
          nama: currentItem.nama,
          jenis: currentItem.jenis,
          qty: parseInt(currentItem.qty)
        });
      } else {
        // Create new item
        await axios.post('/api/atk', {
          nama: currentItem.nama,
          jenis: currentItem.jenis,
          qty: parseInt(currentItem.qty)
        });
      }
      
      setShowModal(false);
      fetchATKItems();
      setError('');
    } catch (err) {
      console.error('Error saving ATK item:', err);
      setError(err.response?.data?.error || 'Gagal menyimpan data');
    } finally {
      setLoading(false);
    }
  };

  // Delete item (soft delete)
  const handleDelete = async (id) => {
    const result = await Swal.fire({
      title: 'Konfirmasi Hapus',
      text: 'Apakah Anda yakin ingin menghapus item ini?',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#ef4444',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'Ya, Hapus!',
      cancelButtonText: 'Batal',
      reverseButtons: true
    });

    if (!result.isConfirmed) {
      return;
    }
    
    setLoading(true);
    try {
      await axios.delete(`/api/atk/${id}`);
      fetchATKItems();
      setError('');
      
      // Show success message
      Swal.fire({
        title: 'Berhasil!',
        text: 'Data ATK berhasil dihapus (soft delete)',
        icon: 'success',
        confirmButtonColor: '#10b981',
        timer: 2000
      });
    } catch (err) {
      console.error('Error deleting ATK item:', err);
      setError('Gagal menghapus data');
      
      // Show error message
      Swal.fire({
        title: 'Gagal!',
        text: err.response?.data?.error || 'Gagal menghapus data',
        icon: 'error',
        confirmButtonColor: '#ef4444'
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Header */}
      <header className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Package className="w-8 h-8 text-indigo-600" />
              <div>
                <h1 className="text-3xl font-bold text-gray-900">Manajemen Inventori ATK</h1>
                <p className="text-sm text-gray-600 mt-1">Sistem Manajemen Alat Tulis Kantor</p>
              </div>
            </div>
            <button
              onClick={handleCreate}
              disabled={loading}
              className="flex items-center space-x-2 bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-lg font-medium transition-all shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Plus className="w-5 h-5" />
              <span>Tambah ATK</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            <p className="font-medium">Error:</p>
            <p className="text-sm">{error}</p>
          </div>
        )}

        {/* Statistics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white p-6 rounded-xl shadow-md border-l-4 border-blue-500">
            <h3 className="text-gray-600 text-sm font-medium">Total Item ATK</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">{atkItems.length}</p>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-md border-l-4 border-green-500">
            <h3 className="text-gray-600 text-sm font-medium">Total Stok</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">
              {atkItems.reduce((sum, item) => sum + item.qty, 0)}
            </p>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-md border-l-4 border-yellow-500">
            <h3 className="text-gray-600 text-sm font-medium">Jenis ATK</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">
              {new Set(atkItems.map(item => item.jenis)).size}
            </p>
          </div>
        </div>

        {/* Table */}
        <div className="bg-white rounded-xl shadow-md overflow-hidden">
          <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900">Daftar Inventori ATK</h2>
          </div>
          
          {loading && !showModal ? (
            <div className="text-center py-12">
              <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-600 border-t-transparent"></div>
              <p className="mt-4 text-gray-600">Memuat data...</p>
            </div>
          ) : atkItems.length === 0 ? (
            <div className="text-center py-12">
              <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <p className="text-gray-500">Belum ada data ATK. Silakan tambah item baru.</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-100 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">ID</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Nama ATK</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Jenis</th>
                    <th className="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">Jumlah Stok</th>
                    <th className="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {atkItems.map((item, index) => (
                    <tr key={item.id} className={`hover:bg-gray-50 transition-colors ${index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}`}>
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{item.id}</td>
                      <td className="px-6 py-4 text-sm text-gray-900 font-medium">{item.nama}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        <span className="px-3 py-1 rounded-full bg-blue-100 text-blue-700 text-xs font-medium">
                          {item.jenis}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm">
                        <span className={`font-semibold ${item.qty < 5 ? 'text-red-600' : item.qty < 20 ? 'text-yellow-600' : 'text-green-600'}`}>
                          {item.qty} pcs
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-center">
                        <div className="flex justify-center space-x-2">
                          <button
                            onClick={() => handleEdit(item)}
                            disabled={loading}
                            className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors disabled:opacity-50"
                            title="Edit"
                          >
                            <Pencil className="w-4 h-4" />
                          </button>
                          <button
                            onClick={() => handleDelete(item.id)}
                            disabled={loading}
                            className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
                            title="Hapus"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </main>

      {/* Modal Form */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-2xl max-w-md w-full">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200">
              <h3 className="text-xl font-semibold text-gray-900">
                {editMode ? 'Edit ATK' : 'Tambah ATK Baru'}
              </h3>
              <button
                onClick={() => setShowModal(false)}
                className="text-gray-400 hover:text-gray-600 transition-colors"
              >
                <X className="w-6 h-6" />
              </button>
            </div>
            
            <form onSubmit={handleSubmit} className="px-6 py-4">
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Nama ATK <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    required
                    value={currentItem.nama}
                    onChange={(e) => setCurrentItem({ ...currentItem, nama: e.target.value })}
                    placeholder="Contoh: Pensil 2B"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition-all"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Jenis <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    required
                    value={currentItem.jenis}
                    onChange={(e) => setCurrentItem({ ...currentItem, jenis: e.target.value })}
                    placeholder="Contoh: Alat Tulis"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition-all"
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Jumlah Stok <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={currentItem.qty}
                    onChange={(e) => setCurrentItem({ ...currentItem, qty: e.target.value })}
                    placeholder="0"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none transition-all"
                  />
                </div>
              </div>
              
              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loading ? 'Menyimpan...' : editMode ? 'Simpan Perubahan' : 'Tambah '}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ATKDashboard;
