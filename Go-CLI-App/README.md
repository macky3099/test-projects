# Booking App
Booking App is a CLI based app built in Go. As per the requirement 3 shows 

## Run Test Cases:
1. go test ./...
2. go test -cover ./...

## How to create executable binaries:
Run the `./build.sh` file on a Linux or MacOS machine to create binaires out of this project.
This binaries can be run simply via the console as per the machine type and its architecture.

# How to test:
1. Run the binary file as per below:
    a. for Windows run it on command prompt by typing `bookingcli.exe` and hitting enter.
    b. for MacOs or Linux run it on bash by typing `./bookingcli-linux-x86` or `./bookingcli-macos-arm` as per the machine and hit enter. 
2. Once you are inside the app then you can select the show number, seats and do the booking.

# Example:
```
manthan-linux@MackPC:/mnt/e/Tech Hunting/TechVerito/bin$ ./bookingcli-linux-x86

Enter Show no: 1
Available Seats:
A1 A2 A3 A4 A5 A6 A7 A8 A9
B1 B2 B3 B4 B5 B6
C1 C2 C3 C4 C5 C6 C7

Enter seats: A1,A2
Print: Sucessfully Booked - Show 1

Subtotal: Rs.640.00
Service Tax @14.00%: Rs.89.60
Swachh Bharat Cess @0.50%: Rs.3.20
Krishi Kalyan Cess @0.50%: Rs.3.20
Total: Rs.736.00

Would you like to continue (Yes/No): Yes


Enter Show no: 1
Available Seats:
A3 A4 A5 A6 A7 A8 A9
B1 B2 B3 B4 B5 B6
C1 C2 C3 C4 C5 C6 C7

Enter seats: A1,A3
Print: A1 NOT available, Please select different seats


Would you like to continue (Yes/No): Yes


Enter Show no: 2
Available Seats:
A1 A2 A3 A4 A5 A6 A7 A8 A9
B1 B2 B3 B4 B5 B6
C1 C2 C3 C4 C5 C6 C7

Enter seats: B1,B2
Print: Sucessfully Booked - Show 2

Subtotal: Rs.560.00
Service Tax @14.00%: Rs.78.40
Swachh Bharat Cess @0.50%: Rs.2.80
Krishi Kalyan Cess @0.50%: Rs.2.80
Total: Rs.644.00

Would you like to continue (Yes/No): No

Total Sales:
Revenue: Rs.1200.00
Service Tax: Rs.168.00
Swachh Bharat Cess: Rs.6.00
Krishi Kalyan Cess: Rs.6.00

manthan-linux@MackPC:/mnt/e/Tech Hunting/TechVerito/bin$
```