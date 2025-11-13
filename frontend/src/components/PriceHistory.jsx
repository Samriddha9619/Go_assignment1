import './PriceHistory.css';

const PriceHistory = ({ history }) => {
  if (!history || history.length === 0) {
    return <div className="no-history">No price history available</div>;
  }

  return (
    <div className="price-history">
      <h4>Price History</h4>
      <div className="history-list">
        {history.map((entry) => (
          <div key={entry.id} className="history-item">
            <span className="history-price">₹{entry.price.toFixed(2)}</span>
            <span className="history-date">
              {new Date(entry.timestamp).toLocaleString()}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default PriceHistory;

