import React, { useEffect, useState, useContext } from "react";
import { UserProfileContext } from '../../App';
import { useParams, useNavigate } from "react-router-dom";
import Navbar from "../NavBar/NavBar";
import Calendar from "../Calendar/Calendar";
import Reservation from "../Reserve/Reserve";
import "./HotelDetails.css"

const HotelDetails = () => {
  const { id } = useParams();
  const [hotel, setHotel] = useState(null);
  const [error, setError] = useState(null);
  const [index, setIndex] = useState(0);
  const [adminError, setAdminError] = useState(null);
  const [amadeusId, setAmadeusId] = useState('');
  const [showAmadeusForm, setShowAmadeusForm] = useState(false)
  const { userProfile } = useContext(UserProfileContext);
  const [selectedDates, setSelectedDates] = useState({
    startDate: new Date(),
    endDate: new Date(),
  });
  const navigate = useNavigate();

  useEffect(() => {
    const fetchHotelDetails = async () => {
      try {
        const response = await fetch(`http://localhost:8085/hotel/${id}`);
        if (response.ok) {
          const data = await response.json();
          setHotel(data);
        } else {
          const errorData = await response.json();
          throw new Error(errorData.error);
        }
      } catch (error) {
        setError(error.message);
      }
    };

    fetchHotelDetails();
  }, [id]);

  const handleSelectDates = (selectedRange) => {
    setSelectedDates(selectedRange);
  };

  const handleDeleteHotel = async () => {
    try {
      const response = await fetch(`http://localhost:8080/hotel/${id}`, {
        method: 'DELETE',
      });
      if (response.ok) {
        navigate(`/`)
      } else {
        const errorData = await response.json();
        throw new Error(errorData.error);
      }
    } catch (error) {
      setAdminError(error.message);
    }
  };

  const handleAmadeusSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch("http://localhost:8090/amadeus", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          hotel_id: id,
          amadeus_id: amadeusId
        }),
      });

      if (response.status === 201) {
        setAmadeusId('');
        setShowAmadeusForm(false);
      } else {
        const errorData = await response.json();
        throw new Error(errorData.error);
      }
    } catch (error) {
      setAdminError(error.message);
    }
  };

  const amadeusForm = showAmadeusForm && (
      <form onSubmit={handleAmadeusSubmit}>
        <input
            type="text"
            placeholder="Ingrese ID de Amadeus"
            value={amadeusId}
            onChange={(e) => setAmadeusId(e.target.value)}
        />
        <button type="submit">Submit</button>
      </form>
  );

  if (error) {
    return (
        <>
          <Navbar />
          <div className="fullscreen">Error: {error}</div>
        </>
    );
  }

  if (!hotel) {
    return (
        <>
          <Navbar />
          <div className="fullscreen">Cargando...</div>
        </>
    );
  }

  // if (!loggedIn) {
  //   return (
  //     <>
  //       <Navbar />
  //       <p className="fullscreen">No puedes acceder a este sitio.</p>
  //     </>
  //   );
  // }

  return (
    <>
      <Navbar />
      <div className="description">
        <h1>{hotel.name}</h1>
        <h4>
          {hotel.street_name} {hotel.street_number}, {hotel.city}
        </h4>
        {hotel.images &&
        <div className="carousel-container">
          <div id={`carousel-${hotel.id}`} className="carousel slide" data-bs-ride="carousel">
            <div className="carousel-inner">
              {hotel.images.map((image) => (
                  <div key={image} className={`carousel-item ${image === hotel.images[index] ? 'active' : ''}`}>
                    <img
                        src={`http://localhost:8080/image?name=${image}`}
                        className="d-block w-100 carousel-img"
                        alt={image.id}
                    />
                  </div>
              ))}
            </div>
            <button className="carousel-control-prev"
                    type="button"
                    data-bs-target={`#carousel-${hotel.id}`}
                    data-bs-slide="prev"
                    onClick={() => setIndex(prevIndex => (prevIndex === 0 ? hotel.images.length - 1 : prevIndex - 1))}>
            <span className="carousel-control-prev-icon" aria-hidden="true"></span>
              <span className="visually-hidden">Previous</span>
            </button>
            <button className="carousel-control-next"
                    type="button"
                    data-bs-target={`#carousel-${hotel.id}`}
                    data-bs-slide="next"
                    onClick={() => setIndex(prevIndex => (prevIndex === hotel.images.length - 1 ? 0 : prevIndex + 1))}>
              <span className="carousel-control-next-icon" aria-hidden="true"></span>
              <span className="visually-hidden">Next</span>
            </button>
          </div>
        </div>
        }

        <p>{hotel.description}</p>
        <h5>Precio por noche: ${hotel.rate}</h5>
        {hotel.amenities && (
            <div>
              <h4>Amenities:</h4>
              <ul className="list">
                {hotel.amenities.map((amenity) => (
                    <li key={amenity}>{amenity}</li>
                ))}
              </ul>
            </div>
        )}

        {userProfile && userProfile.role === "Customer" && (
          <div>
            <h2>Reservar</h2>
            <Calendar onSelectDates={handleSelectDates} />
            <Reservation
              hotel_id={id}
              hotelRate={hotel.rate}
              startDate={selectedDates.startDate}
              endDate={selectedDates.endDate}
            />
          </div>
        )}
        {userProfile && userProfile.role === "Admin" && (
            <div>
              <button className="admin-button" onClick={() => navigate(`/updatehotel/${id}`)}>Modificar Hotel</button>
              <button className="admin-button" onClick={handleDeleteHotel}>Borrar Hotel</button>
              <div>
                {showAmadeusForm ? (
                    amadeusForm
                ) : (
                    <button onClick={() => setShowAmadeusForm(true)}>
                      Agregar ID Amadeus
                    </button>
                )}
              </div>
              {adminError && <p className="error-message">{adminError}</p>}
            </div>
        )}
      </div>
    </>
  );
};

export default HotelDetails;