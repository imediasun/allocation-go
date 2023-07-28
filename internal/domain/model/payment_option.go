package model

import "gitlab.hotel.tools/backend-team/allocation-go/internal/infra/repo/entity"

type PaymentOption string

// Enum values for PaymentOption
const (
	PaymentOptionAcssDebit         PaymentOption = "acss_debit"
	PaymentOptionAfterpayClearpay  PaymentOption = "afterpay_clearpay"
	PaymentOptionAirbnb            PaymentOption = "airbnb"
	PaymentOptionAlipay            PaymentOption = "alipay"
	PaymentOptionAuBecsDebit       PaymentOption = "au_becs_debit"
	PaymentOptionBacsDebit         PaymentOption = "bacs_debit"
	PaymentOptionBancontact        PaymentOption = "bancontact"
	PaymentOptionBankTransfer      PaymentOption = "bank_transfer"
	PaymentOptionBoleto            PaymentOption = "boleto"
	PaymentOptionCash              PaymentOption = "cash"
	PaymentOptionCreditCard        PaymentOption = "credit_card"
	PaymentOptionEps               PaymentOption = "eps"
	PaymentOptionFPX               PaymentOption = "fpx"
	PaymentOptionGiftCard          PaymentOption = "gift_card"
	PaymentOptionGiropay           PaymentOption = "giropay"
	PaymentOptionGrabpay           PaymentOption = "grabpay"
	PaymentOptionIdeal             PaymentOption = "ideal"
	PaymentOptionOnAccount         PaymentOption = "on_account"
	PaymentOptionOxxo              PaymentOption = "oxxo"
	PaymentOptionP24               PaymentOption = "p24"
	PaymentOptionSepaDebit         PaymentOption = "sepa_debit"
	PaymentOptionSofort            PaymentOption = "sofort"
	PaymentOptionSwish             PaymentOption = "swish"
	PaymentOptionVirtualMasterCard PaymentOption = "virtual_master_card"
	PaymentOptionVirtualVisaCard   PaymentOption = "virtual_visa_card"
	PaymentOptionVoucher           PaymentOption = "voucher"
	PaymentOptionWechatPay         PaymentOption = "wechat_pay"
)

func (m *PaymentOption) ToEntity() entity.BookingsNullPaymentOption {
	if m == nil {
		return entity.BookingsNullPaymentOption{}
	}

	return entity.BookingsNullPaymentOptionFrom(entity.BookingsPaymentOption(*m))
}

func PaymentOptionFromEntity(m entity.BookingsNullPaymentOption) *PaymentOption {
	if !m.Valid {
		return nil
	}

	res := PaymentOption(m.Val)
	return &res
}
