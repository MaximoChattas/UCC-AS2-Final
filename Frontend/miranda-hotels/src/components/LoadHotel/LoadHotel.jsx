import React, { useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserProfileContext } from '../../App';
import Navbar from '../NavBar/NavBar';
import './LoadHotel.css';

function LoadHotel() {
    const [name, setName] = useState('');
    const [city, setCity] = useState('');
    const [street_name, setStreet_name] = useState('');
    const [street_number, setStreet_number] = useState('');
    const [room_amount, setRoom_amount] = useState('');
    const [rate, setRate] = useState('');
    const [description, setDescription] = useState('');
    const [amenities, setAmenities] = useState([]);
    const [selectedAmenities, setSelectedAmenities] = useState([]);
    const [images, setImages] = useState([]);

    const [hotelId, setHotelId] = useState('');
    const [isLoaded, setIsLoaded] = useState(false)

    const { userProfile } = useContext(UserProfileContext);

    const [error, setError] = useState('');

    const navigate = useNavigate();

    const handleLoadHotel = async (e) => {
        e.preventDefault();
        setError('');

        try {
            if (!name || !street_name || !street_number || !room_amount || !rate) {
                throw new Error('Complete todos los campos requeridos');
            }

            const response = await fetch('http://localhost:8080/hotel', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name,
                    street_name,
                    city,
                    street_number: parseInt(street_number),
                    room_amount: parseInt(room_amount),
                    rate: parseFloat(rate),
                    description,
                    amenities: selectedAmenities,
                }),
            });

            if (response.status === 201) {
                setIsLoaded(true);
                const data = await response.json();
                setHotelId(data.id);

            } else {
                const data = await response.json();
                const errorMessage = data.error || 'Error';
                throw new Error(errorMessage);
            }
        } catch (error) {
            console.error(error);
            setError(error.message);
        }
    };

    const handleImagesUpload = async (e) => {
        e.preventDefault();

        try {
            const formData = new FormData();
            images.forEach((image) => {
                formData.append('images', image);
            });

            const response = await fetch(`http://localhost:8080/hotel/${hotelId}/images`, {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                navigate('/')
            } else {
                const errorData = await response.json();
                throw new Error(errorData.error);
            }
        } catch (error) {
            console.error(error);
            setError(error.message);
        }
    };

    const handleAmenityChange = (e, amenityName) => {
        if (e.target.checked) {
            setSelectedAmenities([...selectedAmenities, amenityName]);
        } else {
            setSelectedAmenities(selectedAmenities.filter((name) => name !== amenityName));
        }
    };

    useEffect(() => {
        const fetchAmenities = async () => {
            try {
                const response = await fetch('http://localhost:8080/amenity');
                if (response.ok) {
                    const data = await response.json();
                    setAmenities(data);
                } else {
                    const errorData = await response.json();
                    throw new Error(errorData.error);
                }
            } catch (error) {
                console.error(error);
                setError(error.message);
            }
        };

        fetchAmenities();
    }, []);

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files);
        setImages(files);
    };

    if (!userProfile || userProfile.role !== 'Admin') {
        return (
            <>
                <Navbar />
                <p className="contenedorLoad">No puedes acceder a este sitio.</p>
            </>
        );
    }

    return (
        <>
            <Navbar />
            {!isLoaded &&
                <div className="contenedorLoad">
                <h2>Cargar Hotel</h2>
                <form onSubmit={handleLoadHotel}>
                    <div>
                        <label>Nombre:</label>
                        <input type="text" disabled={isLoaded} value={name} onChange={(e) => setName(e.target.value)} />
                    </div>
                    <div>
                        <label>Ciudad:</label>
                        <input type="text" disabled={isLoaded} value={city} onChange={(e) => setCity(e.target.value)} />
                    </div>
                    <div>
                        <label>Calle:</label>
                        <input type="text" disabled={isLoaded} value={street_name} onChange={(e) => setStreet_name(e.target.value)} />
                    </div>
                    <div>
                        <label>Altura:</label>
                        <input
                            type="number"
                            pattern="[0-9]*"
                            value={street_number}
                            disabled={isLoaded}
                            onChange={(e) => setStreet_number(e.target.value)}
                        />
                    </div>
                    <div>
                        <label>Habitaciones:</label>
                        <input
                            type="number"
                            pattern="[0-9]*"
                            value={room_amount}
                            disabled={isLoaded}
                            onChange={(e) => setRoom_amount(e.target.value)}
                        />
                    </div>
                    <div>
                        <label>Tarifa: $</label>
                        <input type="number" disabled={isLoaded} pattern="[0-9]*" value={rate} onChange={(e) => setRate(e.target.value)} />
                    </div>
                    <div>
                        <label>Descripción:</label>
                        <div>
              <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  disabled={isLoaded}
                  placeholder="Ingrese la descripción"
                  maxLength={1000}
                  rows={4}
                  cols={50}
              />
                            <div>
                                Characters disponibles: {1000 - description.length} / {1000}
                            </div>
                        </div>
                    </div>
                    {amenities &&
                    <div>
                        <h5>Amenities:</h5>
                        {amenities.map((amenity) => (
                            <div key={amenity.id}>
                                <input
                                    type="checkbox"
                                    value={amenity.name}
                                    name="amenities"
                                    disabled={isLoaded}
                                    onChange={(e) => handleAmenityChange(e, amenity.name)}
                                />
                                <span>{amenity.name}</span>
                            </div>
                        ))}
                    </div>
                    }
                    {error && <p className="error-message">{error}</p>}
                    <button type="submit" disabled={isLoaded}>
                        {isLoaded ? 'Guardado' : 'Cargar Hotel'}
                    </button>
                </form>
            </div>}

            {isLoaded &&
                <div className="contenedorLoad">
                    <h2>Cargar Imagenes</h2>
                    <form onSubmit={handleImagesUpload}>
                        <input
                            type="file"
                            accept="image/*"
                            multiple
                            disabled={!isLoaded}
                            onChange={handleImageChange}
                        />
                        <button type="submit" disabled={!isLoaded}>
                            Cargar Imagenes
                        </button>
                    </form>
                </div>}
        </>
    );
}

export default LoadHotel;
