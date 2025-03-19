package const_type

type BetStatus string

const BetStatusUnpaid BetStatus = "Unpaid"

const BetStatusPaid BetStatus = "Paid"

const BetStatusInvalid BetStatus = "Invalid"

const BetStatusFailed BetStatus = "Failed"

const BetStatusTimeout BetStatus = "Timeout"

const BetStatusPaying BetStatus = "Paying"

const BetStatusException BetStatus = "Exception"
