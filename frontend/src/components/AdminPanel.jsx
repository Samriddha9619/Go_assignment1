import { useState, useEffect } from 'react';
import { adminAPI } from '../services/api';
import './AdminPanel.css';

const AdminPanel = () => {
  const [logs, setLogs] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    loadLogs();
    loadUsers();
  }, []);

  const loadLogs = async () => {
    try {
      const data = await adminAPI.getScrapingLogs();
      setLogs(data);
    } catch (error) {
      console.error('Error loading logs:', error);
    }
  };

  const loadUsers = async () => {
    try {
      const data = await adminAPI.getAllUsers();
      setUsers(data);
    } catch (error) {
      console.error('Error loading users:', error);
    }
  };

  const handleTriggerScrape = async () => {
    setLoading(true);
    setError('');
    setMessage('');
    try {
      const data = await adminAPI.triggerScrape();
      setMessage(data.message || 'Scraping started successfully');
      setTimeout(() => {
        loadLogs();
      }, 2000);
    } catch (error) {
      setError(error.message || 'Failed to trigger scrape');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h1>Admin Panel</h1>

      <div className="admin-section">
        <h2>Scraping Control</h2>
        <button
          onClick={handleTriggerScrape}
          disabled={loading}
          className="scrape-button"
        >
          {loading ? 'Starting...' : 'Trigger Scrape'}
        </button>
        {message && <div className="success-message">{message}</div>}
        {error && <div className="error-message">{error}</div>}
      </div>

      <div className="admin-section">
        <h2>Scraping Logs</h2>
        <div className="logs-container">
          {logs.length === 0 ? (
            <p>No logs available</p>
          ) : (
            <table className="logs-table">
              <thead>
                <tr>
                  <th>City</th>
                  <th>Status</th>
                  <th>Hotels Count</th>
                  <th>Started At</th>
                  <th>Completed At</th>
                  <th>Error</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => (
                  <tr key={log.id}>
                    <td>{log.city}</td>
                    <td>
                      <span className={`status-badge status-${log.status.toLowerCase()}`}>
                        {log.status}
                      </span>
                    </td>
                    <td>{log.hotels_count}</td>
                    <td>{new Date(log.started_at).toLocaleString()}</td>
                    <td>
                      {log.completed_at
                        ? new Date(log.completed_at).toLocaleString()
                        : '-'}
                    </td>
                    <td className="error-cell">
                      {log.error_message || '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>

      <div className="admin-section">
        <h2>All Users</h2>
        <div className="users-container">
          {users.length === 0 ? (
            <p>No users found</p>
          ) : (
            <table className="users-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Verified</th>
                  <th>Created At</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id}>
                    <td>{user.id}</td>
                    <td>{user.name}</td>
                    <td>{user.email}</td>
                    <td>
                      <span className={`role-badge role-${user.role}`}>
                        {user.role}
                      </span>
                    </td>
                    <td>{user.is_verified ? 'Yes' : 'No'}</td>
                    <td>{new Date(user.created_at).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
};

export default AdminPanel;

