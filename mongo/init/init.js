db = db.getSiblingDB("tournament");

db.createCollection("users");
db.createCollection("tournaments");
db.createCollection("events");
db.createCollection("registrations");
db.createCollection("brackets");
db.createCollection("matches");
db.createCollection("audit_logs");
db.createCollection("appeals");

db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ role: 1 });

db.tournaments.createIndex({ status: 1 });
db.tournaments.createIndex({ created_by: 1 });

db.events.createIndex({ tournament_id: 1 });
db.events.createIndex({ type: 1 });

db.registrations.createIndex({ event_id: 1, player_id: 1 }, { unique: true });
db.registrations.createIndex({ status: 1 });

db.brackets.createIndex({ event_id: 1 }, { unique: true });

db.matches.createIndex({ bracket_id: 1 });
db.matches.createIndex({ event_id: 1, round: 1, position: 1 });
db.matches.createIndex({ referee_id: 1 });

db.audit_logs.createIndex({ operator_id: 1 });
db.audit_logs.createIndex({ target_type: 1, target_id: 1 });
db.audit_logs.createIndex({ created_at: -1 });

db.appeals.createIndex({ appeal_id: 1 }, { unique: true });
db.appeals.createIndex({ match_id: 1 });
db.appeals.createIndex({ appellant_id: 1 });
db.appeals.createIndex({ status: 1 });
db.appeals.createIndex({ created_at: -1 });

db.users.insertOne({
  username: "admin",
  password: "$2a$10$.Hk/3BRk3JcLNJJU/VJiB.TZp3Z7om575eLTJuz39bXoOVzMTDR82",
  display_name: "系统管理员",
  role: "committee",
  created_at: new Date()
});
