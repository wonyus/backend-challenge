// Create application database
db = db.getSiblingDB('backend_challenge');


// Create users collection with indexes
db.createCollection('users');

// Create indexes
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "created_at": 1 });

// Insert sample data (optional)
db.users.insertOne({
  name: "Admin User",
  email: "admin@example.com",
  password: "$2a$06$R.ga34oljt5UqXmSgNR6ze4QpEbq8u9i0Fui/eG2WpZs/nCgjbT1e",
  created_at: new Date(),
  updated_at: new Date()
});

print('Database initialized successfully');