package memcache

import (
	"WB0/pkg/order_struct"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	items             map[int]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[int]Item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC() // данный метод рассматривается ниже
	}
	return &cache
}

func (c *Cache) Set(key int, value interface{}, duration time.Duration) {

	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

	fmt.Println(c.items)

}

func (c *Cache) Get(key int) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	// ключ не найден
	if !found {
		return nil, false
	}

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if item.Expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []int) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []int) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

func (c *Cache) AutoIncrement() int {
	c.RLock()
	defer c.RUnlock()
	id := len(c.items) + 1
	return id
}

func (c *Cache) InsertingDB(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM buy")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Итерация по результатам запроса и заполнение кеша
	for rows.Next() {
		var order order_struct.Order
		var id_order int
		err := rows.Scan(
			&id_order,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSign,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard)
		if err != nil {
			panic(err)
		}

		var delivery order_struct.Delivery

		err = db.QueryRow("SELECT name_receiver, phone_number, zip_code, name_city, address, region, email FROM delivery WHERE delivery_id = $1", id_order).Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			panic(err)
		}
		order.Delivery = delivery

		var payment order_struct.Payment

		err = db.QueryRow("SELECT transaction_id, request_id, name_currency, name_provider, total_amount, payment_dt, name_bank, delivery_cost, goods_total, custom_fee FROM payment WHERE payment_id = $1", id_order).Scan(
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			panic(err)
		}

		order.Payment = payment

		rows_item, err := db.Query("SELECT chrt_id, track_number, price, rid, name_item, sale, size_item, total_price, nm_id, brand, status FROM item WHERE order_id = $1", id_order)
		if err != nil {
			panic(err)
		}
		defer rows_item.Close()

		var itemsOrder []order_struct.Item

		// Итерация по результатам запроса и заполнение кеша
		for rows_item.Next() {
			var item order_struct.Item
			err := rows_item.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status)
			if err != nil {
				panic(err)
			}
			itemsOrder = append(itemsOrder, item)
		}
		order.Items = itemsOrder
		c.Set(id_order, order, 0)
	}

	return nil
}
