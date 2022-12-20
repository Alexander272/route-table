package models

type Analytics struct {
	Number     string `db:"number"`
	OrderStart string `db:"order_start"`
	OrderEnd   string `db:"order_end"`
	Position   int    `db:"position"`
	PosTitle   string `db:"pos_title"`
	Ring       string `db:"ring"`
	PosEnd     string `db:"pos_end"`
	OperTitle  string `db:"oper_title"`
	OperEnd    string `db:"oper_end"`
}
