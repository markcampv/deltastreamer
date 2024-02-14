# DeltaStreamer

## Overview

DeltaStreamer is a tool designed to efficiently monitor and stream deltas (changes) in service states within Consul clusters. It focuses on delivering only the changes in service registration, deregistration, and health status, rather than the entire payload. This approach significantly reduces network traffic and processing overhead, making it ideal for large-scale deployments with frequent changes.

## Features

- **Efficient Change Detection**: Tracks and streams only the deltas in service states, including registrations, deregistrations, and health status updates.
- **Customizable Polling Intervals**: Offers configurable polling intervals to balance between real-time updates and system resource utilization.
- **Scalable Architecture**: Designed to efficiently handle large numbers of services and high rates of change.
- **Simple Integration**: Easy to integrate with existing Consul setups, requiring minimal configuration.

## Getting Started

### Prerequisites

- Go version 1.15 or higher
- Access to a Consul cluster

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/markcampv/deltastreamer.git
