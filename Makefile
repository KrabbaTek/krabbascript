CC      := gcc
CFLAGS  := -Wall -Wextra -fsanitize=address -fno-omit-frame-pointer -static-libasan -g -O2 -I include -I 3rd_party -DKSCRIPT_VENDOR=\"KrabbaTek\" -MMD -MP -std=c11
LDFLAGS := -fsanitize=address -fno-omit-frame-pointer

SRC_DIR  := src
OBJ_DIR  := obj
TEST_DIR := tests
TARGET   := kscript

SRCS    := $(shell find $(SRC_DIR) -name "*.c")
INC     := $(shell find include -name "*.h")
OBJS    := $(SRCS:$(SRC_DIR)/%.c=$(OBJ_DIR)/%.o)
DEPS    := $(OBJS:.o=.d)

LIB_OBJS := $(filter-out $(OBJ_DIR)/main.o, $(OBJS))

TEST_SRCS := $(shell find $(TEST_DIR) -name "*.c")
TEST_BIN  := build/test_runner

TEST_CFLAGS  := -Wall -Wextra -g -O0 -I include -I 3rd_party -DKSCRIPT_VENDOR=\"KrabbaTek\" -DUNITY_INCLUDE_DOUBLE -std=c11 -fsanitize=address -fno-omit-frame-pointer -static-libasan
TEST_LDFLAGS :=

all: $(TARGET)

$(TARGET): $(OBJS)
	@echo "     LD $@"
	@$(CC) $(OBJS) -o $@ $(LDFLAGS)

$(OBJ_DIR)/%.o: $(SRC_DIR)/%.c
	@mkdir -p $(dir $@)
	@echo "     CC $<"
	@$(CC) $(CFLAGS) -c $< -o $@

test: $(TEST_SRCS) $(SRCS)
	@mkdir -p build
	@echo "     LD $(TEST_BIN)"
	@$(CC) $(TEST_CFLAGS) -I$(TEST_DIR) -I3rd_party/unity \
		$(TEST_SRCS) \
		$(filter-out $(SRC_DIR)/main.c, $(SRCS)) \
		3rd_party/unity/unity.c -o $(TEST_BIN)
	@echo "     Running tests..."
	@./$(TEST_BIN)

$(TEST_BIN): $(TEST_SRCS) $(LIB_OBJS)
	@mkdir -p build
	@echo "     LD $@"
	@$(CC) $(TEST_CFLAGS) -I$(TEST_DIR) -I3rd_party/unity $^ 3rd_party/unity/unity.c -o $@ $(TEST_LDFLAGS)

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

.PHONY: all clean install format tidy test

-include $(DEPS)