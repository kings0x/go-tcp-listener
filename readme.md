Creating a tcp connection between a client and server
1. the client connection would have a timeout setting to close if no connection is made after a given time
2. the client would also have a read and write timeout if no packet is sent after some time 
3. the client would have a pinger that would send pings to the listener to show thats its alive
   
4. the listener would have a setDeadline on reads and writes that should timeout if nothing is recieved from the client
5. the listener should have a graceful shutdown mechanism that closes the connection(the client) 