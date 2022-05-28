-to list the services
grpcurl --plaintext localhost:9091 list

-to list all the method
grpcurl --plaintext localhost:9091 list protos.Currency

-to describe the signature of the method
grpcurl --plaintext localhost:9091 describe protos.Currency.GetRate

-to describe the signature
grpcurl --plaintext localhost:9091 describe .protos.RateRequest

-how to curl with data
grpcurl --plaintext -d '{\"Base\": \"IDR\", \"Destination\" : \"USD\"}' localhost:9091 protos.Currency.GetRate