import Navbar from "../NavBar/NavBar.jsx"
import { UserProfileContext } from '../../App';
import React, {useContext, useEffect, useState} from "react";

const ContainerStats = () => {
    const [containers, setContainers] = useState([]);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(true)

    const { userProfile } = useContext(UserProfileContext);

    useEffect(() => {
        const fetchStats = async () => {

            try {
                const response = await fetch("http://localhost:8004/stats");
                if (response.ok) {
                    const data = await response.json();

                    const sortedData = data.sort((a, b) => a.Name.localeCompare(b.Name));
                    setContainers(sortedData);
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

        fetchStats();

        const intervalId = setInterval(fetchStats, 10000);

        return () => clearInterval(intervalId);
        }, [])

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
                <div className="fullscreen">Cargando...</div>
            </>
        )
    }

    if (error) {
        return (
            <>
                <Navbar />
                <div className="fullscreen">Error: {error}</div>
            </>
        );
    }

    if (!containers) {
        return (
            <>
                <Navbar />
                <div className="fullscreen">No hay contenedores corriendo</div>
            </>
        )
    }

    return (
        <>
            <Navbar />
            <h2>Contenedores</h2>
            <table className="container-table fullscreen">
                <thead>
                <tr>
                    <th>ID</th>
                    <th>Nombre</th>
                    <th>CPU (%)</th>
                    <th>Memoria (%)</th>
                    <th>Memoria</th>
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
                    </tr>
                ))}
                </tbody>
            </table>
        </>
    )
}

export default ContainerStats;