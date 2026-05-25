CC      := gcc
CFLAGS  := -Wall -Wextra -fsanitize=address -fno-omit-frame-pointer -static-libasan -g -O2 -I include -DKSCRIPT_VENDOR=\"KrabbaTek\" -MMD -MP -std=c11
LDFLAGS := -fsanitize=address -fno-omit-frame-pointer

SRC_DIR := src
OBJ_DIR := obj
TARGET  := kscript

SRCS    := $(shell find $(SRC_DIR) -name "*.c")
INC     := $(shell find include -name "*.h")
OBJS    := $(SRCS:$(SRC_DIR)/%.c=$(OBJ_DIR)/%.o)
DEPS    := $(OBJS:.o=.d)

all: $(TARGET)

$(TARGET): $(OBJS)
	@echo "     LD $@"
	@$(CC) $(OBJS) -o $@ $(LDFLAGS)

$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c
	@mkdir -p $(dir $@)
	@echo "     CC $<"
	@$(CC) $(CFLAGS) -c $< -o $@

clean:
	@echo "Cleaning up..."
	@rm -rf $(OBJ_DIR) $(TARGET) build

install:
	@echo "Installing..."
	@cp kscript /usr/local/bin/kscript

format:
	@echo "Formatting code..."
	@clang-format -i $(SRCS) $(INC)

tidy:
	@echo "Running clang-tidy..."
	@clang-tidy $(SRCS) -- $(CFLAGS)

.PHONY: all clean install format tidy

-include $(DEPS)