Examine the call stack in a program using runtime.Caller. 

From the golang docs for runtime.Caller

Caller reports file and line number information about function invocations on the calling goroutine's stack. 
The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller. 
(For historical reasons the meaning of skip differs between Caller and Callers.) 
The return values report the program counter, file name, and line number within the file of the corresponding call. 
The boolean ok is false if it was not possible to recover the information.

