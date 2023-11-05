const fs = require('fs');
const MongoClient = require('mongodb').MongoClient;

// Ruta al archivo JSON
const jsonFilePath = './init.json'; // Ajusta la ruta según tu estructura

// URL de conexión a la base de datos MongoDB
const mongoURL = 'mongodb://root:pass@mongodatabase:27017';

// Nombre de la base de datos y colección
const dbName = 'test';
const collectionName = 'hotels';

// Leer el contenido del archivo JSON
fs.readFile(jsonFilePath, 'utf8', (err, data) => {
  if (err) {
    console.error('Error al leer el archivo JSON:', err);
    process.exit(1);
  }

  const jsonData = JSON.parse(data);

  // Conectar a la base de datos MongoDB
  MongoClient.connect(mongoURL, { useUnifiedTopology: true }, (err, client) => {
    if (err) {
      console.error('Error al conectar a MongoDB:', err);
      process.exit(1);
    }

    const db = client.db(dbName);
    const collection = db.collection(collectionName);

    // Insertar los datos del archivo JSON en la colección
    collection.insertMany(jsonData.hotels, (err, result) => {
      if (err) {
        console.error('Error al insertar en MongoDB:', err);
        process.exit(1);
      }

      console.log('Datos insertados en MongoDB:');
      console.log(result.insertedCount + ' documentos insertados');

      // Cerrar la conexión a MongoDB
      client.close();
    });
  });
});