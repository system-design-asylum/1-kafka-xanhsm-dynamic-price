# Problem:  

You are tasked with designing a real-time surge pricing system for a rapidly growing ride-sharing company, similar to Uber or Grab. The goal is to dynamically adjust ride prices based on supply and demand, to incentivize drivers to move to high-demand areas and balance the marketplace.




# Key Requirements:

- Real-time Demand Aggregation: The system needs to ingest real-time data about ride requests from all active users.

- Real-time Supply Aggregation: The system needs to ingest real-time data about available drivers (location, status - e.g., "available," "on trip," "offline").

- Geospatial Analysis: The city should be divided into a grid (or dynamic zones). The system must calculate the demand-to-supply ratio for each grid cell/zone.

- Surge Pricing Calculation: Based on the demand-to-supply ratio in a zone, a surge multiplier (e.g., 1.2x, 1.5x, 2.0x) needs to be calculated. This calculation should be dynamic and potentially configurable.

- Price Application: When a user requests a ride, the system must query the current surge multiplier for their pick-up location and apply it to the base fare before presenting the final price to the user.

- Low Latency: The system must respond with a surge multiplier within 100ms for any given pick-up location.

- Scalability: The system must handle thousands of ride requests and hundreds of thousands of active drivers concurrently. The company plans to expand to many cities globally.

- Fault Tolerance: The system must be highly available and resilient to failures of individual components. Data must not be lost.

- Configurability: The rules for calculating surge (e.g., thresholds for demand-to-supply ratios, maximum surge multipliers, minimum driver density) should be easily configurable without code deployments.

**Non-Functional Considerations (Implied)**:

    Security (basic considerations)

    Cost-effectiveness (general awareness)

    Maintainability




# My Proposed Solution
- Go Services (Producer): A high-performance, concurrent set of microservices written in Golang responsible for:

    Receiving driver heartbeats (longitude, latitude) via an HTTP API.
  
    Receiving customer's requests for ride (pickup coordiation, dropoff coordination, vehicle type) via an HTTP API.

    Performing in-memory point-in-polygon checks against pre-loaded geographic zones to determine each driver's current zone.

    Maintaining in-memory, thread-safe counters for driver supply (and implicitly, ride demand) per zone.

    Periodically batching and updating the drivers table in the database with the latest driver locations.

    Calculating dynamic surge multipliers based on in-memory demand/supply ratios and updating the zones table in batches.

- Apache Kafka: My central nervous system for real-time data streaming.

    Acts as a highly scalable, battle-tested, fault-tolerant buffer for driver heartbeats (though in the current scope, the Go service processes them in-memory before Kafka).

    (Future/Implicit): Can be used to stream aggregated zone metrics or ride requests to other services.

- PostgreSQL with PostGIS: Powerful geospatial database.

    Stores the definitive geographic zone polygons (app_data.zones). These are loaded into the Go service's memory for fast lookups.

    Persists driver locations (app_data.drivers), updated in batched intervals from the Go service.

    Stores the calculated surge multipliers and demand/supply counts for each zone + Strategic in-memory caching for much faster access within high-demands areas.

    Serves spatial read queries from other services (i.e., live map displays) using efficient PostGIS geospatial indexes and calculation functions.


# Principles that I've learnt
- System design is all about making trade-offs all the time.
- Over-the-network is really, really expensive. So perform in-memory processing while you still can.
- It's all about batch processing. Always insert data in batches.


# Technology
Golang

Apache Kafka

PostgreSQL

PostGIS

Docker



# BACKLOG
- [x] Build a baseline system
- [x] Seed small-scale grid data for Ho Chi Minh city map
- [ ] Benchmark dynamic pricing retrieval latency
- [ ] Benchmark components throughput
- [ ] Trace system bottlenecks
- [ ] Fix bottlenecks if any
- [ ] Add load balancers and scale services horizontally; increase request volume to stress test system
- [ ] Handle errors, components outtage
- [ ] Simulate network failure, random corruption within the cluster to see if it withstands
