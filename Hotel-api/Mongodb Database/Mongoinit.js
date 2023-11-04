const fs = require('fs');
const MongoClient = require('mongodb').MongoClient;
const { spawn } = require('child_process');

const dockerComposeFile = 'docker-compose.yml'; // Ruta al archivo Compose de Docker
const mongoURL = 'mongodb://localhost:27017'; // URL de la base de datos MongoDB

const dbName = 'hotel_db'; // Nombre de la base de datos
const collectionName = 'hotels'; // Nombre de la colección

// Conexión a MongoDB
const client = new MongoClient(mongoURL, { useNewUrlParser: true, useUnifiedTopology: true });

async function loadHotels() {
  try {
    const data = fs.readFileSync(dockerComposeFile, 'utf8');
    const composeConfig = JSON.parse(data);
    if (composeConfig.hotels) {
      return composeConfig.hotels;
    }
  } catch (error) {
    console.error(`Error al cargar el archivo Compose de Docker: ${error}`);
  }
  return [];
}

async function saveToMongo(hotels) {
  try {
    const db = client.db(dbName);
    const collection = db.collection(collectionName);

    await collection.deleteMany({});
    const result = await collection.insertMany(hotels);

    console.log(`Se han insertado ${result.insertedCount} registros en MongoDB.`);
  } catch (error) {
    console.error(`Error al insertar en MongoDB: ${error}`);
  }
}

async function watchDockerComposeFile() {
  try {
    let lastModified = fs.statSync(dockerComposeFile).mtimeMs;

    fs.watch(dockerComposeFile, (event, filename) => {
      if (event === 'change') {
        const currentModified = fs.statSync(dockerComposeFile).mtimeMs;
        if (currentModified !== lastModified) {
          console.log(`Se ha detectado un cambio en ${dockerComposeFile}. Cargando hoteles en MongoDB.`);
          lastModified = currentModified;
          loadHotels().then(saveToMongo);
        }
      }
    });
  } catch (error) {
    console.error(`Error al observar el archivo Compose de Docker: ${error}`);
  }
}

(async () => {
  try {
    await client.connect();
    console.log('Conexión a MongoDB establecida.');

    // Cargar los hoteles por primera vez
    const initialHotels = await loadHotels();
    await saveToMongo(initialHotels);

    // Observar cambios en el archivo Docker Compose
    watchDockerComposeFile();
  } catch (error) {
    console.error(`Error: ${error}`);
  }
})();
