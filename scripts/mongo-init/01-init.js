// MongoDB Initialization Script
// This script runs automatically when MongoDB container starts for the first time

// Switch to lumine database
db = db.getSiblingDB('lumine');

// Create collections
db.createCollection('users');
db.createCollection('projects');
db.createCollection('services');

// Create indexes
db.users.createIndex({ "email": 1 }, { unique: true });
db.projects.createIndex({ "name": 1 });
db.services.createIndex({ "type": 1 });

// Insert sample data
db.users.insertOne({
    name: "Admin",
    email: "admin@lumine.local",
    role: "admin",
    createdAt: new Date()
});

print('MongoDB initialization complete!');
