**Problem**:  

You are tasked with designing a real-time surge pricing system for a rapidly growing ride-sharing company, similar to Uber or Grab. The goal is to dynamically adjust ride prices based on supply and demand, to incentivize drivers to move to high-demand areas and balance the marketplace.

**Key Requirements**:

- Real-time Demand Aggregation: The system needs to ingest real-time data about ride requests from all active users.

- Real-time Supply Aggregation: The system needs to ingest real-time data about available drivers (location, status - e.g., "available," "on trip," "offline").

- Geospatial Analysis: The city should be divided into a grid (or dynamic zones). The system must calculate the demand-to-supply ratio for each grid cell/zone.

- Surge Pricing Calculation: Based on the demand-to-supply ratio in a zone, a surge multiplier (e.g., 1.2x, 1.5x, 2.0x) needs to be calculated. This calculation should be dynamic and potentially configurable.

- Price Application: When a user requests a ride, the system must query the current surge multiplier for their pick-up location and apply it to the base fare before presenting the final price to the user.

- Low Latency: The system must respond with a surge multiplier within 100ms for any given pick-up location.

- Scalability: The system must handle thousands of ride requests and hundreds of thousands of active drivers concurrently. The company plans to expand to many cities globally.

- Fault Tolerance: The system must be highly available and resilient to failures of individual components. Data must not be lost.

- Configurability: The rules for calculating surge (e.g., thresholds for demand-to-supply ratios, maximum surge multipliers, minimum driver density) should be easily configurable without code deployments.

Non-Functional Considerations (Implied):

    Security (basic considerations)

    Cost-effectiveness (general awareness)

    Maintainability
