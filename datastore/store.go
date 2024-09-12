package datastore

import (
	nfticketdb "betcube-engine/services/nfticket_service/db"
	stripedb "betcube-engine/services/payment_service/db"
	recorddb "betcube-engine/services/record_service/db"
	sportsbookdb "betcube-engine/services/sportsbook_service/db"
	userdb "betcube-engine/services/user_service/db"
)

type MongoStore struct {
	Sportsbook sportsbookdb.SportsbookStore
}

type SupabaseStore struct {
	User     userdb.UserStore
	Nfticket nfticketdb.NfticketStore
	Record   recorddb.RecordStore
}

type StripeStore struct {
	Payment stripedb.PaymentStore
}
