CC = cc
CFLAGS = -Wall -g
LDFLAGS =

# Source files
SRCS = $(wildcard *.c)
# Object files
OBJS = $(SRCS:.c=.o)
# Executable name
EXEC = sudoku-c

# Default target
all: $(EXEC)

# Link object files to create executable
$(EXEC): $(OBJS)
	$(CC) $(OBJS) -o $(EXEC) $(LDFLAGS)

# Compile source files to object files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Clean up
clean:
	rm -f $(OBJS) $(EXEC)

# Phony targets
.PHONY: all clean