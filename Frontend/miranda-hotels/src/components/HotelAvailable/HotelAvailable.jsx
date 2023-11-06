import React, { useContext, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Navbar from "../NavBar/NavBar";
import "../HotelList/HotelList.css";
import Calendar from "../Calendar/Calendar";
import { format } from "date-fns";

const HotelAvailable = () => {
  const [hotels, setHotels] = useState([]);
  const [error, setError] = useState(null);
  const [city, setCity] = useState('')
  const [selectedDates, setSelectedDates] = useState({
    startDate: new Date(),
    endDate: new Date(),
  });

  const fetchHotels = async () => {
    try {
      const cityFormatted = city.replace(' ', '+')
      const startDate = format(selectedDates.startDate, "dd-MM-yyyy");
      const endDate = format(selectedDates.endDate, "dd-MM-yyyy");
      const startTime = "15:00";
      const endTime = "11:00";
      const startDateTime = `${startDate}+${startTime}`;
      const endDateTime = `${endDate}+${endTime}`;
      const url = `http://localhost:8090/available?city=${cityFormatted}&start_date=${startDateTime}&end_date=${endDateTime}`;
      const response = await fetch(url);
      if (response.ok) {
        const data = await response.json();
        setHotels(data);
      } else {
        const data = await response.json();
        const errorMessage = data.error || "Error";
        throw new Error(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError(error.message);
    }
  };

  const handleSelectDates = (selectedRange) => {
    setSelectedDates(selectedRange);
  };

  if (!hotels) {
    return (
        <>
          <Navbar />
          <h2>Hoteles Disponibles</h2>
          <p className="fullscreen">No hay hoteles disponibles</p>
        </>
    );
  }

  if (hotels.length === 0) {
    return (
      <>
        <Navbar />
        <div className="fullscreen">
          <h2>Verificar Disponibilidad</h2>
          <p>
            Seleccione un rango de fechas en el calendario y una ciudad para verificar los hoteles
            disponibles
          </p>
          <div>
            <label>Ciudad:</label>
            <input type="text" onChange={(e) => setCity(e.target.value)} />
          </div>
          <Calendar onSelectDates={handleSelectDates} />
          <button onClick={fetchHotels} style={{ marginTop: '20px' }}>Verificar</button>
        </div>
      </>
    );
  }

  const formattedStartDate = format(selectedDates.startDate, "dd/MM/yyyy");
  const formattedEndDate = format(selectedDates.endDate, "dd/MM/yyyy");

  return (
    <>
      <Navbar />
      <div className="fullscreen">
        <h2>Hoteles Disponibles</h2>
        <h5>{city}, {formattedStartDate} - {formattedEndDate}</h5>
        <div className="row">
          {hotels.map((hotel) => (
              <div key={hotel.id} className="col-md-4 mb-4">
                <div className="card">
                  {hotel.images &&
                      <img className="card-img-top"
                           alt={`Image for ${hotel.name}`}
                           src={`http://localhost:8080/image?name=${hotel.images[0]}`}
                      />}
                  <div className="card-body">
                    <h5 className="card-title">
                      <Link to={`/hotel/${hotel.id}`}>
                        {hotel.name}
                      </Link>
                    </h5>
                    <p className="card-text">
                      Dirección: {hotel.street_name} {hotel.street_number}
                    </p>
                    <p className="card-text">${hotel.rate}</p>
                  </div>
                </div>
              </div>
          ))}
        </div>
      </div>
    </>
  );
};

export default HotelAvailable;