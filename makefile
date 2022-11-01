namenode:
	bash -c "trap 'go run nameNodeCode/nameNode.go' EXIT";
datanode1:
	bash -c "trap 'go run dataNodeCode/dataNode.go 50051' EXIT";
datanode2:
	bash -c "trap 'go run dataNodeCode/dataNode.go 50052' EXIT";
datanode3:
	bash -c "trap 'go run dataNodeCode/dataNode.go 50053' EXIT";
combine:
	bash -c "trap 'go run combineCode/combine.go' EXIT";
rebeldes:
	bash -c "trap 'go run rebelsCode/rebels.go' EXIT";
clean:
	rm -f dataNodeCode/DATA.txt
	rm -f nameNodeCode/DATA.txt
