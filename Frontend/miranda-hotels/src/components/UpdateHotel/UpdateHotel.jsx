import React, { useContext, useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { UserProfileContext } from '../../App';
import Navbar from '../NavBar/NavBar';
import '../LoadHotel/LoadHotel.css';
import {da} from "react-date-range/dist/locale/index.js";

function UpdateHotel() {
    const { id } = useParams();
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

    const [isLoaded, setIsLoaded] = useState(false)

    const { userProfile } = useContext(UserProfileContext);

    const [error, setError] = useState('');

    const navigate = useNavigate();

    useEffect(() => {
        const fetchHotelDetails = async () => {
            try {
                const response = await fetch(`http://localhost:8080/hotel/${id}`);
                if (response.ok) {
                    const data = await response.json();

                    setName(data.name);
                    setStreet_name(data.street_name);
                    setCity(data.city)
                    setStreet_number(data.street_number.toString());
                    setRoom_amount(data.room_amount.toString());
                    setRate(data.rate.toString());
                    setDescription(data.description);
                    setSelectedAmenities(data.amenities || []);

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

    const handleUpdateHotel = async (e) => {
        e.preventDefault();
        setError('');

        try {
            if (!name || !street_name || !street_number || !room_amount || !rate) {
                throw new Error('Complete todos los campos requeridos');
            }

            const response = await fetch(`http://localhost:8080/hotel/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name,
                    city,
                    street_name,
                    street_number: parseInt(street_number),
                    room_amount: parseInt(room_amount),
                    rate: parseFloat(rate),
                    description,
                    amenities: selectedAmenities,
                }),
            });

            if (response.status === 200) {
                await handleImagesUpload();
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

    const handleAmenityChange = (e, amenityName) => {
        if (e.target.checked) {
            setSelectedAmenities([...selectedAmenities, amenityName]);
        } else {
            setSelectedAmenities(selectedAmenities.filter((name) => name !== amenityName));
        }
    };

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files);
        setImages(files);
    };

    const handleImagesUpload = async () => {
        try {
            const formData = new FormData();
            images.forEach((image) => {
                formData.append('images', image);
            });

            const response = await fetch(`http://localhost:8080/hotel/${id}/images`, {
                method: 'POST',
                body: formData,
            });

            if (response.ok) {
                setTimeout(() => {
                    navigate('/');
                }, 2000);
            } else {
                const errorData = await response.json();
                throw new Error(errorData.error);
            }
        } catch (error) {
            console.error(error);
            setError(error.message);
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
            <div className="contenedorLoad">
                <h2>Modificar Hotel</h2>
                <form onSubmit={handleUpdateHotel}>
                    <div>
                        <label>Nombre:</label>
                        <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
                    </div>
                    <div>
                        <label>Ciudad:</label>
                        <input type="text" value={city} onChange={(e) => setCity(e.target.value)} />
                    </div>
                    <div>
                        <label>Calle:</label>
                        <input type="text" value={street_name} onChange={(e) => setStreet_name(e.target.value)} />
                    </div>
                    <div>
                        <label>Altura:</label>
                        <input
                            type="number"
                            pattern="[0-9]*"
                            value={street_number}
                            onChange={(e) => setStreet_number(e.target.value)}
                        />
                    </div>
                    <div>
                        <label>Habitaciones:</label>
                        <input
                            type="number"
                            pattern="[0-9]*"
                            value={room_amount}
                            onChange={(e) => setRoom_amount(e.target.value)}
                        />
                    </div>
                    <div>
                        <label>Tarifa: $</label>
                        <input type="number" pattern="[0-9]*" value={rate} onChange={(e) => setRate(e.target.value)} />
                    </div>
                    <div>
                        <label>Descripción:</label>
                        <div>
              <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
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
                                        onChange={(e) => handleAmenityChange(e, amenity.name)}
                                        checked={selectedAmenities.includes(amenity.name)}
                                    />
                                    <span>{amenity.name}</span>
                                </div>
                            ))}
                        </div>
                    }

                    <div className="contenedorImagenes">
                        <h4>Cargar Imagenes</h4>
                        <input
                            type="file"
                            accept="image/*"
                            multiple
                            onChange={handleImageChange}
                            disabled={isLoaded}
                        />

                    </div>

                    {error && <p className="error-message">{error}</p>}
                    <button type="submit">
                        Actualizar Hotel
                    </button>
                </form>
            </div>
        </>
    );
}

export default UpdateHotel;
