# Launcher
Thie program allocates a block of memory by calling the system call mmap.
It then executes the target process (passed as the first argument).

This block of memory is later on used by eBPF-based instrumentation to expose values to the target process.

