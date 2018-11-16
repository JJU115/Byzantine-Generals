# Byzantine-Generals
An implementation and demonstration of byzantine fault tolerance in Golang. Specifically the Byzantine Generals OM(m) algorithm.
<br><br>
The algorithm is described in this paper: The Byzantine Generals Problem - Leslie Lamport, Marshall Pease, Robert Shostak - ACM Transactions on Programming Languages and Systems 4, 3 (July 1982), 382-401 
<br>

<h2>Execution</h2>
Invoking the executable requires 3 command line arguments: Recursion level, number of generals, and the commander's order. The recursion level specifies the number of message passing rounds, this must be less than the (number of generals) - 1. The commander's order must be either an "A" (Attack) or "R" (Retreat). <b>Invoke as: BGenerals recursion_level numGenerals A | R</b>

<br>
<br>

<h2>On the implementation</h2>
Generals, or lieutenants, are represented by goroutines. Some generals and even the commander may be traitorous meaning that when they forward messages they wil send the opposite order to even numbered generals. Traitors are chosen at random. <br>
If the number of traitors is no more than a third of the total number of generals then all loyal generals should finish with the same order they originally received from the commander.
