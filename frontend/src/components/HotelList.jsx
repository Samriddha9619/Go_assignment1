import { useState, useEffect } from 'react';
import { hotelAPI } from '../services/api';
import HotelCard from './HotelCard';
import './HotelList.css';

const HotelList = () => {
  const [hotels, setHotels] = useState([]);
  const [cities, setCities] = useState([]);
  const [selectedCity, setSelectedCity] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [sortBy, setSortBy] = useState('price');
  const [sortOrder, setSortOrder] = useState('asc');

  useEffect(() => {
    loadCities();
  }, []);

  useEffect(() => {
    loadHotels();
  }, [selectedCity, sortBy, sortOrder]);

  const loadCities = async () => {
    try {
      const cityList = await hotelAPI.getCities();
      setCities(cityList);
    } catch (error) {
      console.error('Error loading cities:', error);
    }
  };

  const loadHotels = async () => {
    setLoading(true);
    setError('');
    try {
      let data;
      if (selectedCity) {
        data = await hotelAPI.getHotelsByCity(selectedCity, sortBy, sortOrder);
      } else {
        data = await hotelAPI.getAllHotels(sortBy, sortOrder);
      }
      setHotels(data);
    } catch (error) {
      setError(error.message || 'Failed to load hotels');
      console.error('Error loading hotels:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="hotel-list-container">
      <div className="hotel-list-header">
        <h1>Hotels</h1>
        <div className="filters">
          <select
            value={selectedCity}
            onChange={(e) => setSelectedCity(e.target.value)}
            className="city-filter"
          >
            <option value="">All Cities</option>
            {cities.map((city) => (
              <option key={city} value={city}>
                {city}
              </option>
            ))}
          </select>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="sort-filter"
          >
            <option value="price">Sort by Price</option>
            <option value="name">Sort by Name</option>
            <option value="city">Sort by City</option>
          </select>
          <button
            onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
            className="sort-order-button"
          >
            {sortOrder === 'asc' ? '↑' : '↓'}
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      {loading ? (
        <div className="loading">Loading hotels...</div>
      ) : hotels.length === 0 ? (
        <div className="no-hotels">No hotels found</div>
      ) : (
        <div className="hotel-grid">
          {hotels.map((hotel) => (
            <HotelCard key={hotel.id} hotel={hotel} />
          ))}
        </div>
      )}
    </div>
  );
};

export default HotelList;

