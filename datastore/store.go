package datastore

import (
	nfticketdb "betcube_engine/services/nfticket_service/db"
	stripedb "betcube_engine/services/payment_service/db"
	recorddb "betcube_engine/services/record_service/db"
	sportsbookdb "betcube_engine/services/sportsbook_service/db"
	userdb "betcube_engine/services/user_service/db"
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
