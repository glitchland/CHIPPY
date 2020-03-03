###########################################
#
#  S N E K
#
#  Classic game SNEK made for CHIP-8 by 
#  glitch. Shouts to makin-games on NL.
#
#  Press 2/W/Q/E to move the snek 
#
###########################################
: snake_mem_base
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #16
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #32
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #48
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #64	
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #80	
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #96		
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #112		
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #128	
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #144
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #160
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #176
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #192	
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #208	
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #224		
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #240		
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
  0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 #256	

: food_location
  0x00 0x00 # v0 v1

: snake_seg
  0x01 
	
: food_sprite
  0x01
	
: wall_sprite 
  0x01

: padding 
  0x01

# head is allowed to wrap, and will always
# be ahead (hohoho) of the tail - if head
# equals tail at any point the game is over

:alias TO_STORE_X v0
:alias TO_STORE_Y v1
:alias SNAKE_HEAD_X v2
:alias SNAKE_HEAD_Y v3

:alias KEYB_UP v4
:alias KEYB_RIGHT v5
:alias KEYB_LEFT v6
:alias KEYB_DOWN v7
:alias KEYB_EXTEND v8

:alias CURRENT_DIR v9
:alias HEAD_PTR vA
:alias SNAKE_LEN vB

:alias FOOD_Y vC
:alias FOOD_X vD

:alias TEMP_REG vE 
:alias FLAGS_REG vF

:const DIR_UP 0
:const DIR_RIGHT 1
:const DIR_LEFT 2 
:const DIR_DOWN 3

:const FLAG_SET 1 
:const FLAG_UNSET 0

:const SNAKE_SEG_SIZE 2

#
# push the head, increment the head pointer
# 

:const MIN_X 0
:const MAX_X 49
:const MIN_Y 0
:const MAX_Y 31

# sprite vx vy n draw a sprite at x/y position, n rows tall.
# 64x32
: build_walls 
	i := wall_sprite
	#TODO: change 0, and 1 to x and y
	# top
	TO_STORE_X := 0 # use these temporariy
	TO_STORE_Y := 0 # use these temporariy	
	loop
		TO_STORE_X += 1
		while TO_STORE_X < 49
		  sprite TO_STORE_X TO_STORE_Y 1
	again

	# bottom
	TO_STORE_X := 0 # use these temporariy
	TO_STORE_Y := 31 # use these temporariy	
	loop
		TO_STORE_X += 1
		while TO_STORE_X < 49
		  sprite TO_STORE_X TO_STORE_Y 1
	again
	
	# left
	TO_STORE_X := 0 # use these temporariy
	TO_STORE_Y := -1 # use these temporariy	
	loop
		TO_STORE_Y += 1
		while TO_STORE_Y < 32
		  sprite TO_STORE_X TO_STORE_Y 1
	again
	
	# right
	TO_STORE_X := 49 # use these temporariy
	TO_STORE_Y := -1 # use these temporariy	
	loop
		TO_STORE_Y += 1
		while TO_STORE_Y < 32
		  sprite TO_STORE_X TO_STORE_Y 1
	again
	
	TO_STORE_X := 0 
	TO_STORE_Y := 0 
	
	return
	
: push_head_pos_v0_v1
	# save vx save registers v0-vx to i.
	i := snake_mem_base 
	i += HEAD_PTR
	TO_STORE_X := SNAKE_HEAD_X
	TO_STORE_Y := SNAKE_HEAD_Y
	save v1 # save v0, and v1 which is the snake head
	HEAD_PTR += SNAKE_SEG_SIZE
	return
	
: get_tail_pos_v0_v1
	# load vx load registers v0-vx from i.
	# save vx save registers v0-vx to i.
	TEMP_REG := HEAD_PTR
	TEMP_REG -= SNAKE_LEN
	i := snake_mem_base
	i += TEMP_REG
	load v1
	return

: delete_tail_v0_v1
	get_tail_pos_v0_v1		
	i := snake_seg
	sprite v0 v1 1 # delete last segment
	return 

: place_food 
  FOOD_X := random MAX_X
  FOOD_Y := random MAX_Y
	if FOOD_X == 0 then FOOD_X := 1
	if FOOD_X == MAX_X then FOOD_X += -1
	if FOOD_Y == 0 then FOOD_Y := 1
	if FOOD_Y == MAX_Y then FOOD_Y += -1
	
	i := food_sprite
	sprite FOOD_X FOOD_Y 1
	return
	
: check_wall_crash_death
	if SNAKE_HEAD_X < MIN_X then die 
	if SNAKE_HEAD_X > MAX_X then die 
	if SNAKE_HEAD_Y < MIN_Y then die 
	if SNAKE_HEAD_Y > MAX_Y then die
	return

: die
  loop
	again


#
###
: main	
	SNAKE_LEN := SNAKE_SEG_SIZE 
	
	SNAKE_HEAD_Y := 20
  SNAKE_HEAD_X := 20
		
	# keyboard map
	KEYB_UP := 2 # up
	KEYB_RIGHT := 4 # right
	KEYB_LEFT := 6 # left
	KEYB_DOWN := 5 # down (w)
	KEYB_EXTEND := 1 # extend snake
	
	# start the snake off going forward, and initialize
	# the snake state with 2 bytes
	
	TO_STORE_X := SNAKE_HEAD_X
	TO_STORE_Y := SNAKE_HEAD_Y
	push_head_pos_v0_v1
		
	# set the initial direction
	CURRENT_DIR := DIR_RIGHT

	
	# is food eaten?
	# place another food, do sprite outsie of loop
	# handle food
	
	build_walls
	place_food 
		
	# game loop
  loop		
	
		# save current
		push_head_pos_v0_v1
		
		# Set VF to 01 if any set pixels are changed to 
		# unset, and 00 otherwise
		i := snake_seg
		FLAGS_REG := FLAG_UNSET
		sprite SNAKE_HEAD_X SNAKE_HEAD_Y 1
		
		# the snake will overwrite itself if backtrack
		# occurs which sets VF to 1. We can use this to
		# check opposite direction deaths
		if FLAGS_REG == FLAG_SET begin

			TEMP_REG := 0
			FLAGS_REG := FLAG_UNSET
			
			TO_STORE_X := 0x0F
			TO_STORE_X &= SNAKE_HEAD_X
			
			TO_STORE_Y := 0x0F
			TO_STORE_Y &= FOOD_X
			
		  if TO_STORE_X == TO_STORE_Y then TEMP_REG += 1
			
			TO_STORE_X := 0x0F
			TO_STORE_X &= SNAKE_HEAD_Y
			
			TO_STORE_Y := 0x0F
			TO_STORE_Y &= FOOD_Y
			
		  if TO_STORE_X == TO_STORE_Y then TEMP_REG += 1
			
			if TEMP_REG == 2 begin
				#sprite FOOD_Y FOOD_X 1 # erase
				#SNAKE_LEN += SNAKE_SEG_SIZE
				#FLAGS_REG := FLAG_UNSET
				SNAKE_LEN += SNAKE_SEG_SIZE
				push_head_pos_v0_v1
				place_food 
				#jump loop_end
			else
				die
			end
			
		end 
		
		# handle direction input
		if KEYB_UP key then CURRENT_DIR := DIR_UP
		if KEYB_RIGHT key then CURRENT_DIR := DIR_RIGHT 
		if KEYB_LEFT key then CURRENT_DIR := DIR_LEFT
		if KEYB_DOWN key then CURRENT_DIR := DIR_DOWN 

		# adjust sprite coords
		if CURRENT_DIR == DIR_UP then SNAKE_HEAD_Y += -1  
		if CURRENT_DIR == DIR_RIGHT then SNAKE_HEAD_X += -1
		if CURRENT_DIR == DIR_LEFT then SNAKE_HEAD_X += 1
		if CURRENT_DIR == DIR_DOWN then SNAKE_HEAD_Y += 1
				
		check_wall_crash_death
		
		# check edge death
		#if SNAKE_HEAD_X == 0x00 then die
		#if SNAKE_HEAD_X == 0xff then die
		
	  if KEYB_EXTEND key begin
			SNAKE_LEN += SNAKE_SEG_SIZE
		else
			delete_tail_v0_v1	
		end

		# Draw the head first, if the VF flag is
		# set it means it collided. So die.

  again
