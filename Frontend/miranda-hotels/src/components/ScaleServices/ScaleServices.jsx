import Navbar from "../NavBar/NavBar.jsx"
import { UserProfileContext } from '../../App';
import React, {useContext, useEffect, useState} from "react";
import "./ScaleServices.css"
import {el} from "react-date-range/dist/locale/index.js";

const ScaleServices = () => {
    const [options, setOptions] = useState([]);
    const [selectedOption, setSelectedOption] = useState('')
    const [containers, setContainers] = useState([]);
    const [error, setError] = useState(null);
    const [apiError, setApiError] = useState(null)
    const [loading, setLoading] = useState(true)

    const { userProfile } = useContext(UserProfileContext);

    useEffect(() => {
        const fetchOptions = async () => {

            try {
                const response = await fetch("http://localhost:8004/services");
                if (response.ok) {
                    const data = await response.json();
                    setOptions(data);

                    if (data.length > 0) {
                        setSelectedOption(data[0]);

                        const responseStats = await fetch(`http://localhost:8004/stats/${data[0]}`);
                        const stats = await responseStats.json();
                        setContainers(stats);
                    }

                } else {
                    const data = await response.json()
                    const errorMessage = data.error || 'Error';
                    throw new Error(errorMessage);
                }
            } catch (error) {
                console.error(error);
                setError(error.message);
            } finally {
                setLoading(false)
            }
        };

        fetchOptions();

    }, [])

    useEffect(() => {

        const intervalId = setInterval(() => {
            if (selectedOption) {
                fetchStatsData(selectedOption);
            }
        }, 10000);

        return () => clearInterval(intervalId);
    }, [selectedOption]);

    const fetchStatsData = async (option) => {
        try {
            setError(null);
            setApiError(null);
            const response = await fetch(`http://localhost:8004/stats/${option}`);

            if (response.ok) {
                const stats = await response.json();
                setContainers(stats);
            } else {
                const data = await response.json();
                const errorMessage = data.message || 'Error';
                throw new Error(errorMessage)
            }

        } catch (error) {
            console.error('Error fetching stats data:', error);
            setError(error.message);
            await fetchStatsData(selectedOption);
        }
    };

    const handleDropdownChange = async (event) => {
       setLoading(true);
       const selectedValue = event.target.value;
       setSelectedOption(selectedValue);
       await fetchStatsData(selectedValue);
       setLoading(false);
    };

    const handleDeleteContainer = async (containerId) => {
        try {

            const response = await fetch(`http://localhost:8004/container/${containerId}`, {
                method: 'DELETE',
            });

            if (response.ok) {
                const message = await response.text();
                alert(message);
                await fetchStatsData(selectedOption);
            } else {
                const data = await response.json();
                const errorMessage = data.error || 'Error';
                throw new Error(errorMessage);
            }
        } catch (error) {
            setApiError(error.message);
        }
    };

    const handleScaleService = async (option) => {
        try {
            const response = await fetch(`http://localhost:8004/scale/${option}`, {
                method: 'POST',
            });

            if (response.ok) {
                const message = await response.text();
                alert(message);
                await fetchStatsData(option);
            } else {
                const data = await response.json();
                const errorMessage = data.error || 'Error';
                throw new Error(errorMessage);

            }

        } catch (error) {
            console.error('Error scaling service:', error);
            setApiError(error.message);
        }
    }

    if (!userProfile || userProfile.role !== 'Admin') {
        return (
            <>
                <Navbar />
                <p className="fullscreen">No puedes acceder a este sitio.</p>
            </>
        );
    }

    if (loading) {
        return (
            <>
                <Navbar />
                <h2>Contenedores</h2>
                <select id="dropdown" className="dropdown" value={selectedOption} onChange={handleDropdownChange}>
                    {options.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>
                <div className="fullscreen">Cargando...</div>
            </>
        )
    }

    if (error) {
        return (
            <>
                <Navbar />
                <select id="dropdown" className="dropdown" value={selectedOption} onChange={handleDropdownChange}>
                    {options.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>
                <div className="fullscreen">Error: {error}</div>
            </>
        );
    }

    if (!containers) {
        return (
            <>
                <Navbar />
                <select id="dropdown" className="dropdown" value={selectedOption} onChange={handleDropdownChange}>
                    {options.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>
                <div className="fullscreen">No hay contenedores corriendo</div>
            </>
        )
    }

    return (
        <>
            <Navbar />
            <h2>Contenedores</h2>
            <div className="escalate-top-container">
                <select id="dropdown" className="dropdown" value={selectedOption} onChange={handleDropdownChange}>
                    {options.map((option, index) => (
                        <option key={index} value={option}>
                            {option}
                        </option>
                    ))}
                </select>

                <button className="escalate-button" onClick={() => handleScaleService(selectedOption)}>Escalar Servicio</button>
            </div>
            {apiError && <p className="error-message">{apiError}</p>}
            <div className="fullscreen">
                <table className="container-table">
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Nombre</th>
                        <th>CPU (%)</th>
                        <th>Memoria (%)</th>
                        <th>Memoria</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    {containers.map((container) => (
                        <tr key={container.ID}>
                            <td>{container.ID}</td>
                            <td>{container.Name}</td>
                            <td>{container.CPUPerc}</td>
                            <td>{container.MemPerc}</td>
                            <td>{container.MemUsage}</td>
                            <td><button className="escalate-button" onClick={() => {handleDeleteContainer(container.ID)}}>
                                Eliminar
                            </button></td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </>
    )
}

export default ScaleServices;