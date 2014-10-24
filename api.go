package gosberbank

import (
  "encoding/json"
  "fmt"
  "strconv"
)

const (
  registerURL = `https://3dsec.sberbank.ru/payment/rest/register.do`
  reverseURL = `https://3dsec.sberbank.ru/payment/rest/reverse.do`
  refundURL = `https://3dsec.sberbank.ru/payment/rest/refund.do`
  statusURL = `https://3dsec.sberbank.ru/payment/rest/getOrderStatus.do`
  isPaid = 4
      )

//Main API struct
type Sberbank struct {
  currency int
  values map[string][]string
}
// Return initialized *Sberbank
func New(UserName, Password string, Currency int) (s *Sberbank) {
  s = new(Sberbank)
  s.currency = Currency
  s.values = make(map[string][]string)
  s.values["userName"] = []string{UserName}
  s.values["password"] = []string{Password}
  return s
}

type order_js struct {
  OrderId string `json:"orderId"`
  FormUrl string `json:"formUrl"`
  ErrorCode int `json:"errorCode"`
  ErrorMessage string `json:"errorMessage"`
}

//Create new Order
func (s *Sberbank) NewOrder(internalId, amount int, returnUrl string) (err error, order Order) {
  values := s.values
  values["orderNumber"] = []string{strconv.Itoa(internalId)}
  values["amount"] = []string{strconv.Itoa(amount)}
  values["currency"] = []string{strconv.Itoa(s.currency)}
  values["returnUrl"] = []string{returnUrl}
  jsonB, err := sendPost(values, registerURL)
  if err != nil {
    return err, order
  }
  result := new(order_js)
  err = json.Unmarshal(jsonB, result)
  if err != nil {
    return err, order
  }
  if result.ErrorCode != 0 {
    return fmt.Errorf(result.ErrorMessage), order
  }
  order.sberbank = s
  order.SberbankId = result.OrderId
  order.PayPageUrl = result.FormUrl
  order.Amount = amount
  return nil, order
}

//Structure for order
type Order struct {
  //Order Id in sberbank payment
  SberbankId string
  //Url of payment page for this order
  PayPageUrl string
  //Order price
  Amount int
  sberbank *Sberbank
}

type OrderStatus struct {
  OrderStatus int `json:"OrderStatus"`
  ErrorCode int `json:"ErrorCode"`
  ErrorMessage string `json:"ErrorMessage"`
}
// Cancel Order
func (o *Order) Reverse() (err error) {
  values := o.sberbank.values
  values["orderId"] = []string{o.SberbankId}
  jsonB, err := sendPost(values, reverseURL)
  if err != nil {
    return err
  }
  result := new(order_js)
  err = json.Unmarshal(jsonB, result)
  if err != nil {
    return err
  }
  if result.ErrorCode != 0 {
    return fmt.Errorf(result.ErrorMessage)
  }
  return nil
}

//Refund
func (o *Order) Refund() (err error) {
  values := o.sberbank.values
  values["orderId"] = []string{o.SberbankId}
  values["amount"] = []string{strconv.Itoa(o.Amount)}
  jsonB, err := sendPost(values, refundURL)
  if err != nil {
    return err
  }
  result := new(order_js)
  err = json.Unmarshal(jsonB, result)
  if err != nil {
    return err
  }
  if result.ErrorCode != 0 {
    return fmt.Errorf(result.ErrorMessage)
  }
  return nil
}  

//Order Status
func (o *Order) Status() (err error, status *OrderStatus) {
  values := o.sberbank.values
  values["orderId"] = []string{o.SberbankId}
  values["amount"] = []string{strconv.Itoa(o.Amount)}
  jsonB, err := sendPost(values, statusURL)
  if err != nil {
    return err, status
  }
  status = new(OrderStatus)
  err = json.Unmarshal(jsonB, status)
  if err != nil {
    return err, status
  }
  if status.ErrorCode != 0 {
    return fmt.Errorf(status.ErrorMessage), status
  }
  return nil, status
}

//Check paid status
func (o *Order) IsPaid() (paid bool) {
  _, status := o.Status()
  return status.OrderStatus == isPaid
}
