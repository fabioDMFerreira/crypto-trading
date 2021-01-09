rs.reconfig({_id:"rs0",members:[{_id:0,host:"localhost:27017"}]},{force:true})
db.createUser({ user: 'root', pwd: 'password', roles: ['readWrite'] });
