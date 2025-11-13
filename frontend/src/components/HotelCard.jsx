import { useState } from 'react';
import { hotelAPI } from '../services/api';
import PriceHistory from './PriceHistory';
import './HotelCard.css';

const HotelCard = ({ hotel }) => {
  const [showHistory, setShowHistory] = useState(false);
  const [history, setHistory] = useState([]);
  const [loadingHistory, setLoadingHistory] = useState(false);

  const handleViewHistory = async () => {
    if (showHistory) {
      setShowHistory(false);
      return;
    }

    setLoadingHistory(true);
    try {
      const data = await hotelAPI.getPriceHistory(hotel.id);
      setHistory(data);
      setShowHistory(true);
    } catch (error) {
      console.error('Error loading price history:', error);
      alert('Failed to load price history');
    } finally {
      setLoadingHistory(false);
    }
  };

  return (
    <div className="hotel-card">
      <div className="hotel-card-header">
        <h3>{hotel.name}</h3>
        <span className="hotel-city">{hotel.city}</span>
      </div>
      {hotel.location && (
        <p className="hotel-location">📍 {hotel.location}</p>
      )}
      <div className="hotel-price">
        <span className="price-label">Price:</span>
        <span className="price-value">₹{hotel.price.toFixed(2)}</span>
      </div>
      <button
        onClick={handleViewHistory}
        className="history-button"
        disabled={loadingHistory}
      >
        {loadingHistory
          ? 'Loading...'
          : showHistory
          ? 'Hide History'
          : 'View Price History'}
      </button>
      {showHistory && <PriceHistory history={history} />}
    </div>
  );
};

export default HotelCard;

