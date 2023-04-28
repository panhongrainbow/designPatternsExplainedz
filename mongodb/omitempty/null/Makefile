.PHONY: help

no_omit:
	# (1) Alter { Email string `bson:"email,omitempty"` } To { Email string Email string `bson:"email,omitempty"` }
	@sed -i 's/bson:"address,omitempty"/bson:"address"/g' ./omitempty.go
	@echo ""
	# (2) Observe the modified code
	@sed -n '16,24p' ./omitempty.go
	@echo ""
	# (3) Execute omitempty.go
	go run ./omitempty.go
omit:
	# (1) Alter { Email string Email string `bson:"email,omitempty"` } To { Email string `bson:"email,omitempty"` }
	@sed -i 's/bson:"address"/bson:"address,omitempty"/g' ./omitempty.go
	@echo ""
	# (2) Observe the modified code
	@sed -n '16,24p' ./omitempty.go
	@echo ""
	# (3) Execute omitempty.go
	go run ./omitempty.go
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "This makefile has two targets: no_omit and omit."
	@echo "They both modify the omitempty.go file to change the bson tag of the Email field in a struct."
	@echo "Then they run the file and compare the results."
	@echo ""
	@echo "Available targets:"
	@echo "  no_omit  - Modify the address field in the User struct to remove the omitempty option and run the omitempty.go program"
	@echo "  omit     - Modify the address field in the User struct to add the omitempty option and run the omitempty.go program"
	@echo ""