import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import Login from './components/Login';
import Register from './components/Register';
import HotelList from './components/HotelList';
import AdminPanel from './components/AdminPanel';
import ProtectedRoute from './components/ProtectedRoute';
import './App.css';

const AppContent = () => {
  const { isAuthenticated, user, logout, isAdmin } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <Router>
      <div className="app">
        <nav className="navbar">
          <div className="nav-container">
            <Link to="/" className="nav-logo">
              Hotel Scraper
            </Link>
            <div className="nav-links">
              <Link to="/" className="nav-link">
                Hotels
              </Link>
              {isAuthenticated ? (
                <>
                  {isAdmin() && (
                    <Link to="/admin" className="nav-link">
                      Admin
                    </Link>
                  )}
                  <div className="user-info">
                    <span className="user-name">{user?.name}</span>
                    <button onClick={handleLogout} className="logout-button">
                      Logout
                    </button>
                  </div>
                </>
              ) : (
                <Link to="/login" className="nav-link">
                  Login
                </Link>
              )}
            </div>
          </div>
        </nav>

        <main className="main-content">
          <Routes>
            <Route path="/" element={<HotelList />} />
            <Route
              path="/login"
              element={
                isAuthenticated ? (
                  <Navigate to="/" replace />
                ) : (
                  <Login />
                )
              }
            />
            <Route
              path="/register"
              element={
                isAuthenticated ? (
                  <Navigate to="/" replace />
                ) : (
                  <Register />
                )
              }
            />
            <Route
              path="/admin"
              element={
                <ProtectedRoute requireAdmin={true}>
                  <AdminPanel />
                </ProtectedRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
