<h1 align='center'>VANESSA</h1>

![Vanessa_logo](https://user-images.githubusercontent.com/35429211/179373701-b637b68a-8799-42aa-862a-b5d4ce94292a.png)

> VANESSA is a Vehicular Ad-hoc Network (VANET) solution for road safety

## 📝 Table of Contents

<!--ts-->

- [📝 Table of Contents](#-table-of-contents)
- [Problem Statement](#problem-statement)
- [Motivation](#motivation)
- [System Architecture](#system-architecture)
- [Prerequisites](#prerequisites)
- [Environment setup](#environment-setup)
- [Running the project](#running-the-project)
- [How to use](#how-to-use)
  - [The simulator](#the-simulator)
  - [Car UI](#car-ui)
  - [RSU UI](#rsu-ui)
- [Contributing](#contributing)
- [Contributers](#contributers)

<!--te-->

## Problem Statement

<a name="problem-statement"></a>

VANESSA aims to build an efficient and scalable VANET system with protocols that can be used in a wide variety of applications. We used known methods in literature and built on them to meet scalability and efficiency requirements that can be run on any Unix-like system

## Motivation

<a name="motv"></a>

Vehicular Ad-hoc Networks (VANETs) is a rising and active research field that has drawn the attention of many researchers and companies around the globe. The network architecture is continuously and rapidly changing. So, due to frequent disconnection, finding the best routing protocols that can handle such dynamicity is very challenging.

## System Architecture

<a name="sys-arch"></a>

![UI-ARCH](https://user-images.githubusercontent.com/35429211/179374062-feb498c0-3b0d-466b-946d-541ac33f8e9f.png)

## Prerequisites

<a name="prereq"></a>

- Python3 & pip
- Go
- Socat
- Node JS & npm
- A web browser

## Environment setup

<a name="env"></a>

- #### Clone the repo
  ```sh
  $ git clone https://github.com/aaarafat/vanessa.git
  ```
- #### Install mininet wifi

  ```sh
  $ sudo ./scripts/install-mnwifi.sh
  ```

- #### Install dependencies
  ```sh
  $ sudo ./scripts/install-dep.sh
  ```

## Running the project

- #### Build

  ```sh
  $ ./scripts/build.sh
  ```

- #### Running the emulation

  ```sh
  $ sudo ./scripts/run-emulation.sh
  ```

- #### Running the simulation

  ```sh
  $ npm run simulation
  ```

  You can access it on `localhost:4200`

- #### Running the Car UI

  ```sh
  $ npm run car
  ```

  You can access it on `localhost:4201`

- #### Running the RSU UI
  ```sh
  $ npm run rsu
  ```
  You can access it on `localhost:4202`

## How to use

### The simulator

Once we open the simulator, we will see the following
![image](https://user-images.githubusercontent.com/35429211/179374369-57fcbd00-a6b7-4a5d-ba8e-2f76b911833a.png)

We can do the following

- Click on any point at the map to add an initial point
- Export the current state by clicking on Export
- Import an exported file by clicking on Import
- Clear the map by clicking on clear

![image](https://user-images.githubusercontent.com/35429211/179374411-f7fc84a8-fb35-48d0-9758-8bb32438e7dc.png)

we have the following options:

- Add accident on that selected location by clicking on Add Accident on the left-side bar
- Cancel the point by clicking on the rewind button on the top right corner
- Add another point to make a route in order to add a car

![image](https://user-images.githubusercontent.com/35429211/179374451-03dfc3f1-3c2b-48c8-bae4-81b10b7eb612.png)

Now that we have added 2 points and constructed a route
We can add a car or we can still cancel it.

We can adjust the car speed by updating the input shown in the figure below. And then we can add a car by pressing on Add Car button

![image](https://user-images.githubusercontent.com/35429211/179374485-351386c9-c9df-40a9-bf48-cf33ce4a92d9.png)

shows a running simulation containing 2 cars pointed at by 1 and an obstacle pointed at by 2. Also, The RSUs ranges are indicated by the circles pointed at by 3

![image](https://user-images.githubusercontent.com/35429211/179374563-f75f5096-da5f-4dda-9794-461661702231.png)

We can click on any of the added car to show its information and show its current route

![image](https://user-images.githubusercontent.com/35429211/179374584-7e53f162-87ad-494e-a4ac-1168a8282f11.png)

we can see the range of this car and its speed and its port that we will use to login in into its UI.

### Car UI

Once we open the Car UI, we will be asked to enter the port of the car we want to

![image](https://user-images.githubusercontent.com/35429211/179374601-04e8a88d-553c-4758-8ea0-a2744e5222e9.png)

By entering the port that we can get from the simulator we will enter the car UI

![image](https://user-images.githubusercontent.com/35429211/179374614-0f9b0cf7-d18a-4e9f-9ab2-aaa697e1c196.png)

In Car UI, we are shown the current location of the selected car and we can see the messages that it received on the right-side bar.

### RSU UI

Just like the car we will enter the RSU port in the beginning

![image](https://user-images.githubusercontent.com/35429211/179374632-218684b4-b7d6-465f-8123-27dd99fd0f16.png)

And once we enter it, we will be shown the RSU UI, in which we can observe its current state. We can see the Total packets it sent or received

![image](https://user-images.githubusercontent.com/35429211/179374648-13c0e79c-0b14-4ced-890f-8380e3f29432.png)

Also, we can see its ARP table and the list of all obstacles reported by the cars to the RSUs

![image](https://user-images.githubusercontent.com/35429211/179374667-52ba937c-7451-489b-a868-fde0355f0c66.png)

![image](https://user-images.githubusercontent.com/35429211/179374674-af7c73d9-9ce3-4461-ba3f-186cef9ce2f5.png)

![image](https://user-images.githubusercontent.com/35429211/179374680-447f4b5c-4d33-4e75-ad3f-6785fea16bcd.png)

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b AmazingFeature-Feat`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin AmazingFeature-Feat`)
5. Open a Pull Request

## Contributers

<table>
  <tr>
    <td align="center"><a href="https://github.com/D4rk1n"><img src="https://avatars.githubusercontent.com/u/44725090?s=460&v=4" width="100px;" alt=""/><br /><sub><b>Abdelrahman Arafat</b></sub></a><br /></td>
    <td align="center"><a href="https://github.com/fuboki10"><img src="https://avatars.githubusercontent.com/u/35429211?s=460&v=4" width="100px;" alt=""/><br /><sub><b>Abdelrahman Tarek</b></sub></a><br /></td>
    <td align="center"><a href="https://github.com/lido22"><img src="https://avatars.githubusercontent.com/u/42592954?v=4" width="100px;" alt=""/><br /><sub><b>Ahmed Walid</b></sub></a><br /></td>
    <td align="center"><a href="https://github.com/Hassan950"><img src="https://avatars.githubusercontent.com/u/42610032?s=460&v=4" width="100px;" alt=""/><br /><sub><b>Hassan Mohamed</b></sub></a><br /></td>
  </tr>
 </table>
