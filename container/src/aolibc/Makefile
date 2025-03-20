# Builds a patched libc.a, customized with the override functions found in aostdlib.c

# Path to your custom stdlib object file
CUSTOM_LIB_OBJ = aostdio.o

.PHONY: build
build: aolibc.a

aolibc.a: $(CUSTOM_LIB_OBJ)
	$(AR) r aolibc.a $(CUSTOM_LIB_OBJ)

CUSTOM_LIB_OBJ: aostdio.c
	$(CC) -c aostdio.c -o aostdio.o



# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f $(CUSTOM_LIB_OBJ) aolibc.a
	@echo "Clean complete."