// Copyright 2022 Omelchuk Rostyslav <work@rostyslav.io>
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package odb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Opendatabot — платформа для роботи з відкритими державними даними.
// Наше API надає актуальні дані з Єдиного реєстру підприємств та організацій,
// судового реєстру, реєстр виконавчих проваджень, реєстру податкових боргів,
// розкладу судових засідань, реєстру податкових боргів, реєстру платників ПДВ та інших.
//
// Ми вирішуємо задачі з перевірки та контролю контрагентів, моніторингу судових справ
// та судових засідань, перевірки персоналу, отримання нових клієнтів для вашого бізнесу.
//
// Наша мета — надання даних та сервісів для будь-якого програмного забезпечення в Україні.
// Для доступу пишіть нам на mail@opendatabot.com та вкажіть для якого проекту плануєте використання.

const (
	// Компанії та ФОП
	governmentCompaniesEndpoint = "https://opendatabot.com/api/v2/government-companies"
	dpaEndpoint                 = "https://opendatabot.com/api/v2/dpa/%s"
	companyEndpoint             = "https://opendatabot.com/api/v2/company/%s"
	changesEndpoint             = "https://opendatabot.com/api/v2/changed/%s"
	wagedebtEndpoint            = "https://opendatabot.com/api/v2/wagedebt/%s"
	auditEndpoint               = "https://opendatabot.com/api/v2/audit"
	auditByIdEndpoint           = "https://opendatabot.com/api/v2/audit/%s"
	registrationsEndpoint       = "https://opendatabot.com/api/v2/registrations"
	registrationByIdEndpoint    = "https://opendatabot.com/api/v2/registrations/%s"
	inspectionsEndpoint         = "https://opendatabot.com/api/v2/inspections"
	inspectionByIdEndpoint      = "https://opendatabot.com/api/v2/inspections/%s"
	pdfEndpoint                 = "https://opendatabot.com/api/v2/pdf/%s"
	permitsEndpoint             = "https://opendatabot.com/api/v2/permits"
	singletaxEndpoint           = "https://opendatabot.com/api/v2/singletax"
	vatEndpoint                 = "https://opendatabot.com/api/v2/vat"
	// Судовий реєстр
	courtEndpoint               = "https://opendatabot.com/api/v2/court"
	institutionsEndpoint        = "https://opendatabot.com/api/v2/institutions"
	courtByIdEndpoint           = "https://opendatabot.com/api/v2/court/%s"
	scheduleEndpoint            = "https://opendatabot.com/api/v2/schedule"
	accusedEndpoint             = "https://opendatabot.com/api/v2/accused"
	scheduleByIdEndpoint        = "https://opendatabot.com/api/v2/schedule/%s"
	companyCourtsEndpoint       = "https://opendatabot.com/api/v2/company-courts"
	companyCourtsByTypeEndpoint = "https://opendatabot.com/api/v2/company-courts/%s"
	courtCasesEndpoint          = "https://opendatabot.com/api/v2/court-cases/%s"
	// Транспорт
	transportEndpoint             = "https://opendatabot.com/api/v2/transport"
	transportByIdEndpoint         = "https://opendatabot.com/api/v2/transport/%s"
	transportLicensesEndpoint     = "https://opendatabot.com/api/v2/transport-licenses"
	transportLicensesByIdEndpoint = "https://opendatabot.com/api/v2/transport-licenses/%s"
	// Робота API
	genKeyEndpoint     = "https://opendatabot.com/api/v2/genKey"
	statisticsEndpoint = "https://opendatabot.com/api/v2/statistics"
	// Фізичні особи
	alimentEndpoint              = "https://opendatabot.com/api/v2/aliment"
	lawyersEndpoint              = "https://opendatabot.com/api/v2/lawyers"
	lawyersByIdEndpoint          = "https://opendatabot.com/api/v2/lawyers/%s"
	corruptOfficialsByIdEndpoint = "https://opendatabot.com/api/v2/corrupt-officials/%s"
	corruptOfficialsEndpoint     = "https://opendatabot.com/api/v2/corrupt-officials"
	passportEndpoint             = "https://opendatabot.com/api/v2/passport"
	wantedEndpoint               = "https://opendatabot.com/api/v2/wanted"
	// Виконавчі провадження
	fullPenaltyByNumberEndpoint    = "https://opendatabot.com/api/v2/full-penalty/%s"
	fullPenaltyDocByNumberEndpoint = "https://opendatabot.com/api/v2/full-penalty-doc/%s"
	fullPenaltyEndpoint            = "https://opendatabot.com/api/v2/full-penalty"
	performerEndpoint              = "https://opendatabot.com/api/v2/performer"
	penaltiesByCodeEndpoint        = "https://opendatabot.com/api/v2/penalties/%s"
	penaltyByNumberEndpoint        = "https://opendatabot.com/api/v2/penalty/%s"
	penaltiesEndpoint              = "https://opendatabot.com/api/v2/penalties"
	// КОАТУУ
	koatuuRegionsEndpoint       = "https://opendatabot.com/api/v2/koatuu/regions"
	koatuuRegionsByCodeEndpoint = "https://opendatabot.com/api/v2/koatuu/regions/%s"
	// Нерухомість
	realtyEndpoint               = "https://opendatabot.com/api/v2/realty"
	realtyByIdEndpoint           = "https://opendatabot.com/api/v2/realty/%s/%s"
	realtyResultEndpoint         = "https://opendatabot.com/api/v2/realty-result"
	realtyReportByNumberEndpoint = "https://opendatabot.com/api/v2/realty-report/%s"
	// Моніторинг бізнесу
	timelineEndpoint = "https://opendatabot.com/api/v2/timeline"
)

// OdbClient is the main Opendatabot struct of the package
type OdbClient struct {
	Settings *Settings
}

// Option is an option for OdbClient
type Option interface {
	Apply(*Settings)
}

type Settings struct {
	ApiKey string
	Client *http.Client
}

// ApiKey Option
type withApiKey string

func (w withApiKey) Apply(o *Settings) {
	o.ApiKey = string(w)
}

func WithApiKey(apiKey string) Option {
	return withApiKey(apiKey)
}

// NewOdbClient
// Create new client
func NewOdbClient(options ...Option) (*OdbClient, error) {
	settings, err := ApplySettings(options)

	if err != nil {
		return nil, err
	}

	return &OdbClient{Settings: settings}, nil
}

func ApplySettings(options []Option) (*Settings, error) {
	var setting Settings

	for _, option := range options {
		option.Apply(&setting)
	}

	return &setting, nil
}

type GovernmentCompany struct {
	Status string `json:"status"`
	Data   struct {
		Count int `json:"count"`
		Items []struct {
			Code string `json:"code"`
		} `json:"items"`
	} `json:"data"`
}

// GetGovernmentCompany
// Перевірка за кодом, що компанія належить державі
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/government_company
func (odb *OdbClient) GetGovernmentCompany(
	code string, // Код ЄДРПОУ
) (response *GovernmentCompany, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(governmentCompaniesEndpoint, map[string]string{
		"code": code,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "items": [
	//      {
	//        "code": "31325005"
	//      }
	//    ]
	//  }
	//}
}

type FopDpa struct {
	Code     string `json:"code"`      // Код платника податків (ІПН)
	FullName string `json:"full_name"` // ПІБ
	// припинено
	// в стані припинення
	// порушено справу про банкрутство (санація)
	// порушено справу про банкрутство
	// зареєстровано
	// свідоцтво про державну реєстрацію недійсне
	// зареєстровано
	Status             string   `json:"status"`              // Статус
	Phones             []string `json:"phones"`              // Телефони
	Email              string   `json:"email"`               // Електронна пошта
	RegistrationDate   string   `json:"registration_date"`   // Дата реєстрації
	RegistrationNumber string   `json:"registration_number"` // Номер реєстрації
	LastDate           string   `json:"last_date"`           // Дата оновлення інформації
	BirthDate          string   `json:"birth_date"`          // Дата народження
	// male
	// female
	Sex                    string   `json:"sex"`                     // Стать
	Activities             string   `json:"activities"`              // Види діяльності https://tax.gov.ua/yuridichnim-osobam/arhiv/podatki-ta-zbori/ediniy-podatok/perelik-vidiv-diyalnosti---
	AdditionallyActivities []string `json:"additionally_activities"` // Додаткові види діяльності
	ActivityKinds          []struct {
		Name      string `json:"name"` // Назва виду діяльності
		Code      string `json:"code"` // Код виду діяльності
		IsPrimary bool   `json:"is_primary"`
	} `json:"activity_kinds"`
	Registrations []struct {
		EndDate     string `json:"end_date"`    // Дата зняття з обліку
		Code        string `json:"code"`        // Ідентифікаційний код органу
		Name        string `json:"name"`        // Назва органу
		Description string `json:"description"` // Опис взяття на облік
		Type        string `json:"type"`        // Тип взяття на облік
		StartDate   string `json:"start_date"`  // Дата взяття на облік
	} `json:"registrations"` // Дані реєстраторів
	Registration struct {
		Date         string `json:"date"`          // Дата реєстрації
		RecordNumber string `json:"record_number"` // Номер реєстрації
		RecordDate   string `json:"record_date"`   // Дата запису
	} `json:"registration"`
	Termination struct {
		State        int    `json:"state"`
		StateText    string `json:"state_text"`
		Date         string `json:"date"`
		RecordNumber string `json:"record_number"` // Номер реєстрації
		Cause        string `json:"cause"`         // Причина припинення
	} `json:"termination"` // Статус припинення
	TerminationCancel struct {
		Date         string `json:"date"`
		RecordNumber string `json:"record_number"` // Номер реєстрації
		DocDate      string `json:"doc_date"`
		CourtName    string `json:"court_name"`
		DocNumber    string `json:"doc_number"`
		DateJudge    string `json:"date_judge"`
	} `json:"termination_cancel"`
	History []struct {
		Date    string `json:"date"` // Дата внесення змін
		Changes []struct {
			// full_name
			// ceo_name
			// location
			Field    string `json:"field"`     // Поле в якому відбулися зміни
			OldValue string `json:"old_value"` // Старе значення
			NewValue string `json:"new_value"` // Нове значення
		} `json:"changes"` // Зміни
	} `json:"history"` // Історія змін
	PdvCode   string `json:"pdv_code"`   // Код ПДВ
	PdvStatus string `json:"pdv_status"` // Статус ПДВ
	TaxDebts  struct {
		Text         string `json:"text"`          // Текстова інформація
		Icon         string `json:"icon"`          // ⚠️
		Total        string `json:"total"`         // Загальний податковий борг
		Local        string `json:"local"`         // Місцевий податковий борг
		Government   string `json:"government"`    // Державний податковий борг
		DatabaseDate string `json:"database_date"` // Дата актуальності
		Type         string `json:"type"`          //
	} `json:"tax_debts"`
	Singletax struct {
		DateStart string `json:"date_start"` // Дата відкриття єдиного податку
		DateEnd   string `json:"date_end"`   // Дата закриття єдиного податку
		Rate      string `json:"rate"`       // Відсоткова ставка єдиного податку
		Group     string `json:"group"`      // Група податку
		Active    bool   `json:"active"`     // Статус єдиного податку
	} `json:"singletax"`
	SingletaxRisk struct {
		Text string `json:"text"` // Текстова інформація
		Icon string `json:"icon"` // ⚠️
	} `json:"singletax_risk"` // Можлива втрата Єдиного податку
	Address struct {
		Zip     string `json:"zip"`
		Country string `json:"country"`
		Address string `json:"address"` // Адреса
		Parts   struct {
			Atu       string `json:"atu"`
			AtuCode   string `json:"atu_code"`
			Street    string `json:"street"`
			HouseType string `json:"house_type"`
			House     string `json:"house"`
			Building  string `json:"building"`
			NumType   string `json:"num_type"`
			Num       string `json:"num"`
		} `json:"parts"`
	} `json:"address"` // Блок з адресою
	TaxDepartments struct {
		TaxDepartmentId         int    `json:"tax_department_id"`
		CREG                    int    `json:"C_REG"`
		CDST                    int    `json:"C_DST"`
		CRAJ                    int    `json:"C_RAJ"`
		NAMERAJ                 string `json:"NAME_RAJ"`
		TSTI                    int    `json:"T_STI"`
		NAMESTI                 string `json:"NAME_STI"`
		CSTI                    int    `json:"C_STI"`
		Code                    int    `json:"code"`
		KoatuuCode              string `json:"koatuu_code"`
		RegionTaxDepartmentCode int    `json:"region_tax_department_code"`
	} `json:"tax_departments"`
	TaxRequisites []struct {
		Type      string `json:"type"`
		KoatuuObl string `json:"koatuu_obl"`
		Koatuu    string `json:"koatuu"`
		Location  string `json:"location"`
		Recipient string `json:"recipient"`
		Code      int    `json:"code"`
		Bank      string `json:"bank"`
		Mfo       int    `json:"mfo"`
		Iban      string `json:"iban"`
		TaxCode   int    `json:"tax_code"`
	} `json:"tax_requisites"`
}

// GetDpa Отримання реєстраційної інформації ФОП
// (ПІБ, адреса, види діяльності, статус, дані з ДПА)
// за індівідуальним кодом платника податків (ІПН), статус платника ПДВ
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/dpa
func (odb *OdbClient) GetDpa(
	code string, // індівідуальний код платника податків (ІПН)
) (response *FopDpa, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(dpaEndpoint, code)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "code": "1111111111",
	//  "full_name": "Петров Іван Володимирович",
	//  "status": "зареєстровано",
	//  "phones": [
	//    "+380111111111"
	//  ],
	//  "email": "mail@email.com",
	//  "registration_date": "2017-01-01",
	//  "registration_number": "11220000000034567",
	//  "last_date": "2018-01-01 12:37:14",
	//  "birth_date": "1987-11-03",
	//  "sex": "male",
	//  "activities": "62.01 Комп'ютерне програмування",
	//  "additionally_activities": [
	//    "62.01 Комп'ютерне програмування"
	//  ],
	//  "activity_kinds": [
	//    {
	//      "name": "Видання комп'ютерних ігор",
	//      "code": "58.21",
	//      "is_primary": true
	//    }
	//  ],
	//  "registrations": [
	//    {
	//      "end_date": "2017-04-18",
	//      "code": "39484073",
	//      "name": "КРОПИВНИЦЬКА ОБ'ЄДНАНА ДЕРЖАВНА ПОДАТКОВА IНСПЕКЦIЯ ГОЛОВНОГО УПРАВЛIННЯ ДФС У КIРОВОГРАДСЬКIЙ ОБЛАСТI",
	//      "description": "дані про взяття на облік як платника податків",
	//      "type": "taxoffice",
	//      "start_date": "2017-04-18"
	//    }
	//  ],
	//  "registration": {
	//    "date": "2017-01-01",
	//    "record_number": "12340123300054405",
	//    "record_date": "2008-03-19"
	//  },
	//  "termination": {
	//    "state": 3,
	//    "state_text": "припинено",
	//    "date": "2019-11-01",
	//    "record_number": "11220000000034567",
	//    "cause": "Припинення ФОП за її рішенням"
	//  },
	//  "termination_cancel": {
	//    "date": "2019-11-01",
	//    "record_number": "11220000000034567",
	//    "doc_date": "2019-10-01",
	//    "court_name": "постанова Вищий господарський суд України",
	//    "doc_number": "5016/3089/2012(18/8)",
	//    "date_judge": "2019-10-01"
	//  },
	//  "history": [
	//    {
	//      "date": "2017-04-18",
	//      "changes": [
	//        {
	//          "field": "ceo_name",
	//          "old_value": "Петрова Галина Сергіївна",
	//          "new_value": "Шевченко Галина Сергіївна"
	//        }
	//      ]
	//    }
	//  ],
	//  "pdv_code": "143605704021",
	//  "pdv_status": "active",
	//  "tax_debts": {
	//    "text": "Податковий борг на 01.02.2018 — 37 334 000 грн",
	//    "icon": "⚠️",
	//    "total": "37334",
	//    "local": "3861.00",
	//    "government": "33473.00",
	//    "database_date": "01.02.2018",
	//    "type": "available"
	//  },
	//  "singletax": {
	//    "date_start": "2016-01-01",
	//    "date_end": "2016-01-01",
	//    "rate": "5",
	//    "group": "3",
	//    "active": true
	//  },
	//  "singletax_risk": {
	//    "text": "Можлива втрата Єдиного податку, податковий борг більше 6 місяців",
	//    "icon": "⚠️"
	//  },
	//  "address": {
	//    "zip": "49000",
	//    "country": "УКРАЇНА",
	//    "address": "01034, м.Київ, Шевченківський район, ВУЛИЦЯ ЯРОСЛАВІВ ВАЛ, будинок 55, корпус Б",
	//    "parts": {
	//      "atu": "Дніпропетровська обл., місто Дніпро, Жовтневий район",
	//      "atu_code": "1222166908",
	//      "street": "ВУЛИЦЯ ПОЛЯ",
	//      "house_type": "буд.",
	//      "house": "2",
	//      "building": "3",
	//      "num_type": "кв.",
	//      "num": "402"
	//    }
	//  },
	//  "tax_departments": {
	//    "tax_department_id": 463,
	//    "C_REG": 4,
	//    "C_DST": 63,
	//    "C_RAJ": 63,
	//    "NAME_RAJ": "СОБОРНИЙ РАЙОН М.ДНIПРА",
	//    "T_STI": 63,
	//    "NAME_STI": "ГУ ДФС У ДНIПРОПЕТРОВСЬКIЙ ОБЛ.(СОБОРНИЙ Р-Н М.ДНIПРА)",
	//    "C_STI": 39394856,
	//    "code": 0,
	//    "koatuu_code": "1210136900",
	//    "region_tax_department_code": 43145015
	//  },
	//  "tax_requisites": [
	//    {
	//      "type": "singletax",
	//      "koatuu_obl": "1200000000",
	//      "koatuu": "1210136900",
	//      "location": "СОБОРНИЙ",
	//      "recipient": "УК у Собор.р.м.Дніпра/Собор.р/18050401",
	//      "code": 37989269,
	//      "bank": "Казначейство України(ел. адм. подат.)",
	//      "mfo": 899998,
	//      "iban": "UA788999980334139866000004005",
	//      "tax_code": 18050401
	//    }
	//  ]
	//}
}

type CompanyData struct {
	FullName      string `json:"full_name"`  // Повна назва компанії
	ShortName     string `json:"short_name"` // Скорочена назва компанії
	Code          string `json:"code"`       // Код ЄДРПОУ
	CeoName       string `json:"ceo_name"`   // ПІБ
	Location      string `json:"location"`   // Адреса
	Activities    string `json:"activities"` // Види діяльності
	Status        string `json:"status"`     // зареєстровано, зареєстровано, свідоцтво про державну реєстрацію недійсне, порушено справу про банкрутство, порушено справу про банкрутство (санація), в стані припинення, припинено
	Beneficiaries []struct {
		Title    string `json:"title"`    // ПІБ
		Capital  int64  `json:"capital"`  // Капітал
		Location string `json:"location"` // Адреса
	} `json:"beneficiaries"`
	DatabaseDate string `json:"database_date"` // Дата оновлення інформації
	PdvCode      string `json:"pdv_code"`      // Код ПДВ
	PdvStatus    string `json:"pdv_status"`    // Статус ПДВ
}

// GetCompany
// Отримання реєстраційної інформації за кодом ЄДРПОУ.
// Можливо отримати запит за декількома компаніями —
// коди повинні передаватись через кому без пробілів.
// Наприклад: 41711425,32746583,39896792
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/Company
func (odb *OdbClient) GetCompany(
	code string, // коди ЄДРПОУ
) (response []CompanyData, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(companyEndpoint, code)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//[
	//  {
	//    "full_name": "ПУБЛІЧНЕ АКЦІОНЕРНЕ ТОВАРИСТВО КОМЕРЦІЙНИЙ БАНК 'ПРИВАТБАНК'",
	//    "short_name": "ПАТ КБ 'ПРИВАТБАНК'",
	//    "code": "11111111",
	//    "ceo_name": "Петров Іван Володимирович",
	//    "location": "01034, м.Київ, Шевченківський район, ВУЛИЦЯ ЯРОСЛАВІВ ВАЛ, будинок 55, корпус Б",
	//    "activities": "62.01 Комп'ютерне програмування",
	//    "status": "зареєстровано",
	//    "beneficiaries": [
	//      {
	//        "title": "Петров Іван Володимирович",
	//        "capital": 206059743960,
	//        "location": "01034, м.Київ, Шевченківський район, ВУЛИЦЯ ЯРОСЛАВІВ ВАЛ, будинок 55, корпус Б"
	//      }
	//    ],
	//    "database_date": "2018-01-01 19:04:32",
	//    "pdv_code": "143605704021",
	//    "pdv_status": "active"
	//  }
	//]
}

type ChangeData struct {
	Code  string `json:"code"` // Код ЄДРПОУ
	Items []struct {
		Date    string `json:"date"` // Дата внесення змін
		Changes []struct {
			Field    string `json:"field"`     // Поле в якому відбулися зміни
			OldValue string `json:"old_value"` // Старе значення
			NewValue string `json:"new_value"` // Нове значення
		} `json:"changes"` // Зміни
	} `json:"items"`
}

// GetChanges
// Отримання переліку змін реєстраційної інформації
// (зміна директора, адреси, статусу, види діяльності, назви, власників).
// Для даних за декількома компаніями — коди повинні передаватись через кому без пробілів.
// Наприклад: 41711425,32746583,39896792.
// Параметр from передається у форматі YYYY-MM-DD
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/changes
func (odb *OdbClient) GetChanges(
	code string, // коди ЄДРПОУ
	params map[string]string, //map[string]string{
	//	"from":	"дата, з якої показати зміни",
	//}
) (response []ChangeData, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(changesEndpoint, code)

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//[
	//  {
	//    "code": "11111111",
	//    "items": [
	//      {
	//        "date": "2017-04-18",
	//        "changes": [
	//          {
	//            "field": "ceo_name",
	//            "old_value": "Петрова Галина Сергіївна",
	//            "new_value": "Шевченко Галина Сергіївна"
	//          }
	//        ]
	//      }
	//    ]
	//  }
	//]
}

type Wagedebt struct {
	Code           string `json:"code"`            // Код ЄДРПОУ
	Debt           string `json:"debt"`            // Сумма заборгованості
	PenaltiesCount string `json:"penalties_count"` // Кількість виконавчіх проваджень
	Name           string `json:"name"`            // Повна назва компанії
	DatabaseDate   string `json:"database_date"`   // Дата актуальності
	Active         int    `json:"active"`          // Ознака актуальності
}

// GetWagedebt
// Отримання публічної інформації щодо компанії,
// наявність в базі боржників по заробітній платі
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/wagedebt
func (odb *OdbClient) GetWagedebt(
	code string, // код ЄДРПОУ
) (response *Wagedebt, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(wagedebtEndpoint, code)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "code": "11111111",
	//  "debt": "13682.43",
	//  "penalties_count": "2️",
	//  "name": "ПУБЛІЧНЕ АКЦІОНЕРНЕ ТОВАРИСТВО КОМЕРЦІЙНИЙ БАНК 'ПРИВАТБАНК'",
	//  "database_date": "2018-05-25",
	//  "active": 1
	//}
}

type AuditsData struct {
	AuditId string `json:"audit_id"` // Внутрішній id
	Code    string `json:"code"`     // Код компанії/внутрішній id ФОПа
	Date    string `json:"date"`     // Дата перевірки
	Type    string `json:"type"`     // Вид перевірки
	Pib     string `json:"pib"`      // Ім'я ФОП
}

// GetAudit
// Отримання публічної інформації щодо проведення планових перевірок
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/audits
func (odb *OdbClient) GetAudit(
	params map[string]string, //map[string]string{
	//	"code":		"код ОКПО",
	//	"pib":		"Ім'я ФОП",
	//	"limit":	"Кількість записів",
	//	"offset":	"Зміщення",
	//}
) (response []AuditsData, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(auditEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//[
	//  {
	//    "audit_id": "13682",
	//    "code": "13682.43",
	//    "date": "2018-05-25",
	//    "type": "DEBT",
	//    "pib": "РОМАНІВ МИКОЛА ІВАНОВИЧ"
	//  }
	//]
}

// GetAuditById
// Отримання публічної інформації щодо проведення планових перевірок
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/audit
func (odb *OdbClient) GetAuditById(
	id string, // внутрішній id
) (response []AuditsData, err error) {
	if err = checkNotEmpty(id); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(auditByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//[
	//  {
	//    "audit_id": "13682",
	//    "code": "13682.43",
	//    "date": "2018-05-25",
	//    "type": "DEBT",
	//    "pib": "РОМАНІВ МИКОЛА ІВАНОВИЧ"
	//  }
	//]
}

type Registrations struct {
	Count int `json:"count"` // Кількість збігів
	Items []struct {
		Id               string `json:"id"`                // ідентифікатор запису
		Type             string `json:"type"`              // Тип юридична (1) або фізична (2) особа
		FullName         string `json:"full_name"`         // Повна назва компанії
		Activity         string `json:"activity"`          // Види діяльності
		RegistrationDate string `json:"registration_date"` // Дата реєстрації
		RegionId         int    `json:"region_id"`         // ідентифікатор регіону
	} `json:"items"`
}

// GetRegistrations
// Отримання переліку нових компаній та ФОПів
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/registrations
func (odb *OdbClient) GetRegistrations(
	params map[string]string, //map[string]string{
	//	"offset": 			"Зміщення відносно початку результатів пошуку",
	//	"limit": 			"Кількість записів",
	//	"type": 			"юридична (company) або фізична (fop) особа",
	//	"reg_date_from":	"пошук за датою з YYYY-MM-DD",
	//	"reg_date_to": 		"пошук за датою по YYYY-MM-DD",
	//	"activities": 		"сортування за видами діяльності, через OR, наприклад, 69 OR 96",
	//	"location": 		"пошук за адресою, Дніпро OR київ",
	//	"is_phone": 		"Фільтр по наявності телефону [0|1]",
	//	"is_email": 		"Фільтр по наявності email [0|1]",
	//	"sort": 			"спосіб сортувааня (за зростанням 'ASC' або спаданням'DESC')",
	//}
) (response *Registrations, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(registrationsEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "count": 1,
	//  "items": [
	//    {
	//      "id": "273987",
	//      "type": "1",
	//      "full_name": "ПУБЛІЧНЕ АКЦІОНЕРНЕ ТОВАРИСТВО КОМЕРЦІЙНИЙ БАНК 'ПРИВАТБАНК'",
	//      "activity": "62.01 Комп'ютерне програмування",
	//      "registration_date": "2017-01-01",
	//      "region_id": 1
	//    }
	//  ]
	//}
}

type Registration struct {
	Code      string `json:"code"`
	FullName  string `json:"full_name"`  // Повна назва компанії
	ShortName string `json:"short_name"` // Скорочена назва компанії
	Location  string `json:"location"`   // Адреса
	CeoName   string `json:"ceo_name"`   // ПІБ
	Activity  string `json:"activity"`   // Види діяльності
	Status    string `json:"status"`     // Статус
	// зареєстровано
	// зареєстровано, свідоцтво про державну реєстрацію недійсне
	// порушено справу про банкрутство
	// порушено справу про банкрутство (санація)
	// в стані припинення, припинено
	Email            string `json:"email"`             // Електронна пошта
	Phones           string `json:"phones"`            // Телефони
	RegistrationDate string `json:"registration_date"` // Дата реєстрації
	Capital          string `json:"capital"`           // Капітал
	Type             string `json:"type"`              // Тип юридична (1) або фізична (2) особа
	RegionId         int    `json:"region_id"`         // Iдентифікатор регіону
}

// GetRegistrationById
// Отримання реєстраційної інформації за внутрішнім id
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/registration
func (odb *OdbClient) GetRegistrationById(
	id string, // внутрішній id, який отримали з пошуку нових компаній/ФОПів
) (response *Registration, err error) {
	if err = checkNotEmpty(id); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(registrationByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "code": "11111111",
	//  "full_name": "ПУБЛІЧНЕ АКЦІОНЕРНЕ ТОВАРИСТВО КОМЕРЦІЙНИЙ БАНК 'ПРИВАТБАНК'",
	//  "short_name": "ПАТ КБ 'ПРИВАТБАНК'",
	//  "location": "01034, м.Київ, Шевченківський район, ВУЛИЦЯ ЯРОСЛАВІВ ВАЛ, будинок 55, корпус Б",
	//  "ceo_name": "Петров Іван Володимирович",
	//  "activity": "62.01 Комп'ютерне програмування",
	//  "status": "зареєстровано",
	//  "email": "mail@email.com",
	//  "phones": "+380111111111,+380222222222",
	//  "registration_date": "2017-01-01",
	//  "capital": "59743960",
	//  "type": "1",
	//  "region_id": 1
	//}
}

type InspectionsResponse struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість перевірок
		Items []struct {
			Id              string `json:"id"`               // ідентифікатор запису
			Code            string `json:"code"`             // код ЄДРПОУ
			Name            string `json:"name"`             // Перевіряючий орган
			Address         string `json:"address"`          // Адреса
			Region          string `json:"region"`           // Ідентифікатор регіону
			Status          string `json:"status"`           // Статус перевірки
			Risk            string `json:"risk"`             // Ризик
			LastModify      string `json:"last_modify"`      // Час останьої модифікації
			DateStart       string `json:"date_start"`       // Час початку перевірки
			DateEnd         string `json:"date_end"`         // Час закінчення перевірки
			Regulator       string `json:"regulator"`        // Регулятор
			ParentRegulator string `json:"parent_regulator"` // Головне управління регулятора
			ActivityType    string `json:"activity_type"`    // Ціль перевірки
			DatabaseDate    string `json:"database_date"`    // Дата додання у базу
			ViolationsCount string `json:"violations_count"` // Кількість порушень
			PartsCount      string `json:"parts_count"`      // Кількість результатів переврок
		} `json:"items"`
	} `json:"data"`
}

// GetInspections
// Отримання інформації про перевірки
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/inspections
func (odb *OdbClient) GetInspections(
	code string, // код ЄДРПОУ
) (response *InspectionsResponse, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(inspectionsEndpoint, map[string]string{
		"code": code,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "id": "2570834",
	//        "code": "43145015",
	//        "name": "ГОЛОВНЕ УПРАВЛІННЯ ДПС У ДНІПРОПЕТРОВСЬКІЙ ОБЛАСТІ",
	//        "address": "країна, 49005, Дніпропетровська обл., місто Дніпро, ВУЛИЦЯ СІМФЕРОПОЛЬСЬКА, будинок 17-А",
	//        "region": "4",
	//        "status": "Проведено",
	//        "risk": "Незначний",
	//        "last_modify": "2021-05-07 09:28:00",
	//        "date_start": "2021-05-07 09:28:00",
	//        "date_end": "2021-05-07 09:28:00",
	//        "regulator": "Головне управління Пенсійного фонду України у Дніпропетровській області",
	//        "parent_regulator": "Пенсійний фонд України",
	//        "activity_type": "дотримання суб'єктом господарювання вимог законодавства у сфері загальнообов'язкового державного пенсійного страхування щодо достовірності поданих відомостей",
	//        "database_date": "2021-05-08",
	//        "violations_count": "1",
	//        "parts_count": "1"
	//      }
	//    ]
	//  }
	//}
}

type InspectionItemResponse struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Id              string `json:"id"`               // ідентифікатор запису
		Code            string `json:"code"`             // код ЄДРПОУ
		Name            string `json:"name"`             // Перевіряючий орган
		Address         string `json:"address"`          // Адреса
		Region          string `json:"region"`           // Ідентифікатор регіону
		Status          string `json:"status"`           // Статус перевірки
		Risk            string `json:"risk"`             // Ризик
		LastModify      string `json:"last_modify"`      // Час останьої модифікації
		DateStart       string `json:"date_start"`       // Час початку перевірки
		DateEnd         string `json:"date_end"`         // Час закінчення перевірки
		Regulator       string `json:"regulator"`        // Регулятор
		ParentRegulator string `json:"parent_regulator"` // Головне управління регулятора
		ActivityType    string `json:"activity_type"`    // Ціль перевірки
		DatabaseDate    string `json:"database_date"`    // Дата додання у базу
		ViolationsCount string `json:"violations_count"` // Кількість порушень
		PartsCount      string `json:"parts_count"`      // Кількість результатів переврок
	} `json:"data"`
}

// GetInspectionById
// Отримання інформації про перевірку за ідентифікатором перевірки
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/inspection-id
func (odb *OdbClient) GetInspectionById(
	id string, // Ідентифікатор перевірки
) (response *InspectionItemResponse, err error) {
	if err = checkNotEmpty(id); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(inspectionByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "id": "2570834",
	//    "code": "43145015",
	//    "name": "ГОЛОВНЕ УПРАВЛІННЯ ДПС У ДНІПРОПЕТРОВСЬКІЙ ОБЛАСТІ",
	//    "address": "країна, 49005, Дніпропетровська обл., місто Дніпро, ВУЛИЦЯ СІМФЕРОПОЛЬСЬКА, будинок 17-А",
	//    "region": "4",
	//    "status": "Проведено",
	//    "risk": "Незначний",
	//    "last_modify": "2021-05-07 09:28:00",
	//    "date_start": "2021-05-07 09:28:00",
	//    "date_end": "2021-05-07 09:28:00",
	//    "regulator": "Головне управління Пенсійного фонду України у Дніпропетровській області",
	//    "parent_regulator": "Пенсійний фонд України",
	//    "activity_type": "дотримання суб'єктом господарювання вимог законодавства у сфері загальнообов'язкового державного пенсійного страхування щодо достовірності поданих відомостей",
	//    "database_date": "2021-05-08",
	//    "violations_count": "1",
	//    "parts_count": "1"
	//  }
	//}
}

type Pdf struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Link string `json:"link"` // Посилання на документ
	} `json:"data"`
}

// GetPdf
// Генерування pdf документа з повною інформацією за кодом ЄДРПОУ
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/pdf
func (odb *OdbClient) GetPdf(
	code string, // код ЄДРПОУ
) (response *Pdf, err error) {
	if err = checkNotEmpty(code); err != nil {
		return nil, err
	}

	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(pdfEndpoint, code)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "link": "https://opendatabot.com/pdf/inn/91/1234567890-91908-39710-96d5c444edaак3cd79084ab59b25ed24.pdf"
	//  }
	//}
}

type LicensesData struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		Count string `json:"count"` // Кількість знайдених об'єктів
		Items []struct {
			Number string `json:"number"` // Registration number
			Type   string `json:"type"`
			// Пальне
			// Спирт
			// Виробництво пального
			// Зберігання пального
			// Оптова торгівля пальним, за відсутності місць оптової торгівлі
			// Оптова торгівля пальним, за наявності місць оптової торгівлі
			// Роздрібна торгівля пальним
			// Зберігання пального (виключно для потреб власного споживання чи промислової переробки)
			Subtype          string `json:"subtype"`
			StartDate        string `json:"start_date,omitempty"`
			EndDate          string `json:"end_date,omitempty"`
			RenewalDate      string `json:"renewal_date,omitempty"`
			PauseDate        string `json:"pause_date,omitempty"`
			CancelationDate  string `json:"cancelation_date,omitempty"`
			Active           int    `json:"active"`
			Address          string `json:"address,omitempty"`
			RegistrationDate string `json:"registration_date,omitempty"`
		} `json:"items"`
	} `json:"data"`
}

// GetPermits
// Отримати інформацію щодо ліцензій компанії
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/licensesData
func (odb *OdbClient) GetPermits(
	params map[string]string, //map[string]string{
	//	"code":	"код ЄДРПОУ або ІПН",
	//	"pib":	"Статус ліцензії. Available values : 0, 1",
	//}
) (response *LicensesData, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(permitsEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "items": [
	//      {
	//        "number": "123123",
	//        "type": "oil_license",
	//        "subtype": "Пальне",
	//        "start_date": "2017-01-01",
	//        "end_date": "2017-01-01",
	//        "renewal_date": "2017-01-01",
	//        "pause_date": "2017-01-01",
	//        "cancelation_date": "2017-01-01",
	//        "active": 1
	//      },
	//      {
	//        "number": "321312",
	//        "type": "oil_excise",
	//        "subtype": "Пальне",
	//        "address": "string",
	//        "registration_date": "2017-01-01",
	//        "active": 1
	//      }
	//    ]
	//  }
	//}
}

type SingletaxSuccess struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		Count string `json:"count"` // Кількість знайдених об'єктів
		Items []struct {
			FopHash   string `json:"fop_hash"`
			Name      string `json:"name"`       // Назва компанії
			Code      string `json:"code"`       // Код компанії
			DateStart string `json:"date_start"` // Дата відкриття єдиного податку
			DateEnd   string `json:"date_end"`   // Дата закриття єдиного податку
			Rate      string `json:"rate"`       // Відсоткова ставка єдиного податку
			Group     string `json:"group"`      // Група податку
			Active    bool   `json:"active"`     // Статус єдиного податку
		} `json:"items"`
	} `json:"data"`
}

// GetSingletax
// Отримати інформацію щодо єдиного податку компанії
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/singletax
func (odb *OdbClient) GetSingletax(
	params map[string]string, //map[string]string{
	//	"code": 	"код ЄДРПОУ або ІПН",
	//	"pib": 		"ПІБ людини",
	//	"fophash": 	"Хеш фізичної особи",
	//}
) (response *SingletaxSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(singletaxEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "items": [
	//      {
	//        "fop_hash": "sdkjb2372ryfwnd",
	//        "name": "МЦ АСКЛЕПІЙ",
	//        "code": "45035236",
	//        "date_start": "2016-01-01",
	//        "date_end": "2016-01-01",
	//        "rate": "5",
	//        "group": "3",
	//        "active": true
	//      }
	//    ]
	//  }
	//}
}

type Vat struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		PdvCode      string `json:"pdv_code"`      // Код ПДВ
		PdvStatus    string `json:"pdv_status"`    // Статус платника
		DateAnul     string `json:"date_anul"`     // Дата анулювання> (якщо анульовано)
		Name         string `json:"name"`          // Повна назва компанії або ПІБ ФОП
		Code         string `json:"code"`          // Код компанії (якщо компанія)
		DatabaseDate string `json:"database_date"` // Дата оновлення інформації
	} `json:"data"`
}

// GetVat
// Отримання інформації по коду платника ПДВ
// https://docs.opendatabot.com/#/%D0%9A%D0%BE%D0%BC%D0%BF%D0%B0%D0%BD%D1%96%D1%97%20%D1%82%D0%B0%20%D0%A4%D0%9E%D0%9F/vat
func (odb *OdbClient) GetVat(
	params map[string]string, //map[string]string{
	//	"vatNumber": 	"Код ПДВ",
	//	"ipn": 			"Код ІПН",
	//	"companyCode":	"Код компанії",
	//}
) (response *Vat, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(vatEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "pdv_code": "1234567890",
	//    "pdv_status": "nonactive",
	//    "date_anul": "2017-01-18",
	//    "name": "Петров Іван Володимирович",
	//    "code": "1234567890",
	//    "database_date": "2018-04-02"
	//  }
	//}
}

type CourtDecisions struct {
	Status string `json:"status"` // Статус операції
	Count  int    `json:"count"`  // Кількість збігів
	Items  []struct {
		DocId        int    `json:"doc_id"`        // Внутрішній id
		CourtCode    int    `json:"court_code"`    // Внутрішній код судової установи
		CourtName    string `json:"court_name"`    // Назва судової установи
		JudgmentCode int    `json:"judgment_code"` // Внутрішній код Форми судочинства
		// Кримінальне
		// Цивільне
		// Господарське
		// Адміністративне
		// Адмінправопорушення
		JudgmentName string `json:"judgment_name"` // Форма судочинства
		JusticeCode  int    `json:"justice_code"`  // Внутрішній код Типу процесуального документа
		// Вирок
		// Постанова
		// Рішення
		// Судовий наказ
		// Ухвала
		// Окрема ухвала
		// Окрема думка
		JusticeName      string `json:"justice_name"`      // Тип процесуального документа
		CategoryCode     int    `json:"category_code"`     // Внутрішній код категорії справи
		CategoryName     string `json:"category_name"`     // Категорія справи
		CauseNumber      string `json:"cause_number"`      // Номер справи
		AdjudicationDate string `json:"adjudication_date"` // Дата набрання законної сили
		DatePubl         string `json:"date_publ"`         // Дата публікації
		ReceiptDate      string `json:"receipt_date"`      // Дата реєстрації
		Judge            string `json:"judge"`             // Суддя
		Link             string `json:"link"`              // Посилання на рішення
	} `json:"items"`
}

// GetCourt
// Отримання судових рішень
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/courtDecision
func (odb *OdbClient) GetCourt(
	params map[string]string, //map[string]string{
	//	// 1: Цивільне
	//	// 2: Кримінальне
	//	// 3: Господарське
	//	// 4: Адміністративне
	//	// 5: Адмінправопорушення
	//	"judgment_code": "1",
	//	// 1: Вирок
	//	// 2: Постанова
	//	// 3: Рішення
	//	// 4: Судовий наказ
	//	// 5: Ухвала
	//	// 6: Окрема ухвала
	//	// 10: Окрема думка
	//	"justice_code": "1",
	//	"court_code":   "Код суда (перелік в судових реєстрах по /institutions)",
	//	"company_code": "код ЄДРПОУ компанії",
	//	"text":         "Пошук в тексті рішення",
	//	// first
	//	// appeal
	//	// cassation
	//	"stage":           "Тип інстанциї",
	//	"text_intro":      "Пошук в вступній частині рішення",
	//	"text_resolution": "Пошук в резолютивній частині рішення",
	//	"offset":          "Зміщення відносно початку результатів пошуку",
	//	"limit":           "Кількість записів",
	//	"date_from":       "Зміщення від дати ухвали рішення",
	//	"date_to":         "Зміщення до дати ухвали рішення",
	//	"number":          "Номер справи",
	//	"search_criteria": "Критерій пошуку значення параметру text в тексті судового рішення. words_in_a_row - Слова повинні йти один за одним",
	//}
) (response *CourtDecisions, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(courtEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "count": 1,
	//  "items": [
	//    {
	//      "doc_id": 1,
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judgment_code": 3,
	//      "judgment_name": "Кримінальне",
	//      "justice_code": 2,
	//      "justice_name": "Вирок",
	//      "category_code": 2,
	//      "category_name": "Клопотання слідчого, прокурора, сторони кримінального провадження",
	//      "cause_number": "509/3997/18",
	//      "adjudication_date": "2018-09-03 00:00:00+03",
	//      "date_publ": "2018-09-03 00:00:00+03",
	//      "receipt_date": "2018-09-03 00:00:00+03",
	//      "judge": "Попов І.В.",
	//      "link": "https://opendatabot.com/court/76195512-b22b9fbfbb31d6aca70b89d1257287a4"
	//    }
	//  ]
	//}
}

type Institution struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count string `json:"count"` // Кількість знайдених судів
		Items []struct {
			Name     string `json:"name"`      // Найменування суду
			CourtId  string `json:"court_id"`  // ID судової установі
			Code     string `json:"code"`      // Код суду
			RegionId string `json:"region_id"` // Номер регіону
			Stage    string `json:"stage"`     // Інстанція
			TypeId   string `json:"type_id"`   // Тип суду
		} `json:"items"`
	} `json:"data"`
}

// GetInstitutions
// Отримання судів
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/institutionsDictionary
func (odb *OdbClient) GetInstitutions(
	params map[string]string, //map[string]string{
	//	"name": 	"Найменування суду",
	//	"offset": 	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//}
) (response *Institution, err error) {
	err = odb.Do(institutionsEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "items": [
	//      {
	//        "name": "Амур-Нижньодніпровський районний суд м.Дніпропетровська",
	//        "court_id": "66",
	//        "code": "0401",
	//        "region_id": "4",
	//        "stage": "first",
	//        "type_id": "3"
	//      }
	//    ]
	//  }
	//}
}

type CourtItem struct {
	DocId            int    `json:"doc_id"`            // Внутрішній id
	CourtCode        int    `json:"court_code"`        // Внутрішній код судової установи
	CourtName        string `json:"court_name"`        // Назва судової установи
	JudgmentCode     int    `json:"judgment_code"`     // Внутрішній код Форми судочинства
	JudgmentName     string `json:"judgment_name"`     // Форма судочинства
	JusticeCode      int    `json:"justice_code"`      // Внутрішній код Типу процесуального документа
	JusticeName      string `json:"justice_name"`      // Тип процесуального документа
	CategoryCode     int    `json:"category_code"`     // Внутрішній код категорії справи
	CategoryName     string `json:"category_name"`     // Категорія справи
	CauseNumber      string `json:"cause_number"`      // Номер справи
	AdjudicationDate string `json:"adjudication_date"` // Дата набрання законної сили
	DatePubl         string `json:"date_publ"`         // Дата публікації
	ReceiptDate      string `json:"receipt_date"`      // Дата реєстрації
	Judge            string `json:"judge"`             // Суддя
	DocumentLink     string `json:"document_link"`     // Посилання на докумен
	Text             string `json:"text"`              // Текст документа
}

// GetCourtById
// Отримання судового документа
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/courtItem
func (odb *OdbClient) GetCourtById(
	id string, // id судового документа
) (response *CourtItem, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(courtByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	// {
	//  "doc_id": 1,
	//  "court_code": 1521,
	//  "court_name": "Овідіопольський районний суд Одеської області",
	//  "judgment_code": 3,
	//  "judgment_name": "Кримінальне|Цивільне|Господарське|Адміністративне|Адмінправопорушення",
	//  "justice_code": 2,
	//  "justice_name": "Вирок|Постанова|Рішення|Судовий наказ|Ухвала|Окрема ухвала|Окрема думка",
	//  "category_code": 2,
	//  "category_name": "Клопотання слідчого, прокурора, сторони кримінального провадження",
	//  "cause_number": "509/3997/18",
	//  "adjudication_date": "2018-09-03 00:00:00+03",
	//  "date_publ": "2018-09-03 00:00:00+03",
	//  "receipt_date": "2018-09-03 00:00:00+03",
	//  "judge": "Попов І.В.",
	//  "document_link": "https://opendatabot.ua/court/81716443-0e094a02739d9ec66a4192c55eec9117",
	//  "text": "string"
	//}
}

type Schedule struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			HearingId    string   `json:"hearing_id"`    // ID судової справи
			Judge        string   `json:"judge"`         // Піб судді
			Forma        string   `json:"forma"`         // Форма судочинства
			Number       string   `json:"number"`        // Номер справи
			CourtId      string   `json:"court_id"`      // id судової установи
			Involved     string   `json:"involved"`      // Позивач/відповідач
			Description  string   `json:"description"`   // Опис справи
			Date         string   `json:"date"`          // Дата та час засідання
			JudgmentCode string   `json:"judgment_code"` // внутрішній код судочинства
			Code         string   `json:"code"`          // Код суду
			Accused      []string `json:"accused"`       // Список звинувачених
		} `json:"items"`
	} `json:"data"`
}

// GetSchedule
// Пошук в судовому розкладі
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/schedule
func (odb *OdbClient) GetSchedule(
	params map[string]string, //map[string]string{
	//	"text_involved":    "Пошук в тексті",
	//	"text_description": "Пошук в описі",
	//	"date":             "Пошук по даті",
	//	"courtId":          "Пошук по courtId",
	//	"offset":           "Зміщення відносно початку результатів пошуку",
	//	"limit":            "Кількість записів",
	//	"judgment_code":    "Внутрішній код Форми судочинства",
	//	"number":           "Пошук по номеру справи",
	//	"date_from":        "Фільтр за датою події (Y-m-d)",
	//	"date_to":          "Фільтр за датою події (Y-m-d)",
	//	"region_id":        "Ідентифікатор регіону", //1 - Автономна Республіка Крим
	//	//2 - Вінницька обл
	//	//3 - Волинська обл
	//	//4 - Дніпропетровська обл
	//	//5 - Донецька обл
	//	//6 - Житомирська обл
	//	//7 - Закарпатська обл
	//	//8 - Запорізька обл
	//	//9 - Івано-Франківська обл
	//	//10 - Київська обл
	//	//11 - Кіровоградська обл
	//	//12 - Луганська обл
	//	//13 - Львівська обл
	//	//14 - Миколаївська обл
	//	//15 - Одеська обл
	//	//16 - Полтавська обл
	//	//17 - Рівненська обл
	//	//18 - Сумська обл
	//	//19 - Тернопільська обл
	//	//20 - Харківська обл
	//	//21 - Херсонська обл
	//	//22 - Хмельницька обл
	//	//23 - Черкаська обл
	//	//24 - Чернівецька обл
	//	//25 - Чернігівська обл
	//	//26 - м.Київ
	//	//27 - м.Севастополь
	//}
) (response *Schedule, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(scheduleEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "hearing_id": "ID судової справи",
	//        "judge": "Піб судді",
	//        "forma": "Адміністративні справи",
	//        "number": "823/1768/17",
	//        "court_id": "690",
	//        "involved": "Позивач (заявник): Головне управління ДФС у Черкаській області, відповідач (боржник): Товариство з обмеженою відповідальністю БМБ Маргарин",
	//        "description": "про стягнення коштів з рахунків у банках",
	//        "date": "2017-11-27 09:30:00",
	//        "judgment_code": "4",
	//        "code": "0401",
	//        "accused": [
	//          "Метельський Андрій Володимирович"
	//        ]
	//      }
	//    ]
	//  }
	//}
}

type Accused struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Forma        string   `json:"forma"`         // Форма судочинства
			Number       string   `json:"number"`        // номер справи
			CourtId      string   `json:"court_id"`      // id судової установи
			Description  string   `json:"description"`   // Опис справи
			JudgmentCode string   `json:"judgment_code"` // внутрішній код судочинства
			Accused      []string `json:"accused"`       // Список звинувачених
		} `json:"items"`
	} `json:"data"`
}

// GetAccused
// Пошук осіб, які обвинувачюються у вчиненні кримінальних та адміністративних правопорушень
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/accused
func (odb *OdbClient) GetAccused(
	params map[string]string, //map[string]string{
	//	"offset":			"Зміщення відносно початку результатів пошуку",
	//	"limit":			"Кількість записів",
	//	"judgment_code":	"Внутрішній код Форми судочинства",
	//	"article":			"Стаття Кримінального кодексу або Кодексу про адміністративні правопорушення",
	//	"region_id":		"Ідентифікатор регіону", //1 - Автономна Республіка Крим
	//	//2 - Вінницька обл
	//	//3 - Волинська обл
	//	//4 - Дніпропетровська обл
	//	//5 - Донецька обл
	//	//6 - Житомирська обл
	//	//7 - Закарпатська обл
	//	//8 - Запорізька обл
	//	//9 - Івано-Франківська обл
	//	//10 - Київська обл
	//	//11 - Кіровоградська обл
	//	//12 - Луганська обл
	//	//13 - Львівська обл
	//	//14 - Миколаївська обл
	//	//15 - Одеська обл
	//	//16 - Полтавська обл
	//	//17 - Рівненська обл
	//	//18 - Сумська обл
	//	//19 - Тернопільська обл
	//	//20 - Харківська обл
	//	//21 - Херсонська обл
	//	//22 - Хмельницька обл
	//	//23 - Черкаська обл
	//	//24 - Чернівецька обл
	//	//25 - Чернігівська обл
	//	//26 - м.Київ
	//	//27 - м.Севастополь
	//	"pib":				"ПІБ обвинуваченного або правопорушника",
	//	"date_from":		"Початкова дата пошуку (Y-m-d)",
	//	"date_to":			"Кінцева дата пошуку (Y-m-d)",
	//}
) (response *Accused, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(accusedEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "forma": "Адміністративні справи",
	//        "number": "823/1768/17",
	//        "court_id": "690",
	//        "description": "про стягнення коштів з рахунків у банках",
	//        "judgment_code": "4",
	//        "accused": [
	//          "Метельський Андрій Володимирович"
	//        ]
	//      }
	//    ]
	//  }
	//}
}

type ScheduleItemMain struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		HearingId    string   `json:"hearing_id"`    // ID судової справи
		Judge        string   `json:"judge"`         // Піб судді
		Forma        string   `json:"forma"`         // Форма судочинства
		Number       string   `json:"number"`        // Номер справи
		CourtId      string   `json:"court_id"`      // id судової установи
		Involved     string   `json:"involved"`      // Позивач/відповідач
		Description  string   `json:"description"`   // Опис справи
		Date         string   `json:"date"`          // Дата та час засідання
		JudgmentCode string   `json:"judgment_code"` // внутрішній код судочинства
		Code         string   `json:"code"`          // Код суду
		Accused      []string `json:"accused"`       // Список звинувачених
	} `json:"data"`
}

// GetScheduleById
// Пошук судового засідання по ID
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/schedule_by_id
func (odb *OdbClient) GetScheduleById(
	id string, // ID судового засідання
) (response *ScheduleItemMain, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(scheduleByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "hearing_id": "ID судової справи",
	//    "judge": "Піб судді",
	//    "forma": "Адміністративні справи",
	//    "number": "823/1768/17",
	//    "court_id": "690",
	//    "involved": "Позивач (заявник): Головне управління ДФС у Черкаській області, відповідач (боржник): Товариство з обмеженою відповідальністю БМБ Маргарин",
	//    "description": "про стягнення коштів з рахунків у банках",
	//    "date": "2017-11-27 09:30:00",
	//    "judgment_code": "4",
	//    "code": "0401",
	//    "accused": [
	//      "Метельський Андрій Володимирович"
	//    ]
	//  }
	//}
}

type CompanyCourtsList struct {
	Civil struct {
		Count     string `json:"count"`      // Кількість виконавчіх проваджень
		LiveCount string `json:"live_count"` // Кількість справ за якими заплановані засідання
	} `json:"civil"`
	Criminal struct {
		Count     string `json:"count"`      // Кількість виконавчіх проваджень
		LiveCount string `json:"live_count"` // Кількість справ за якими заплановані засідання
	} `json:"criminal"`
	Arbitrage struct {
		Count     string `json:"count"`      // Кількість виконавчіх проваджень
		LiveCount string `json:"live_count"` // Кількість справ за якими заплановані засідання
	} `json:"arbitrage"`
	Administrative struct {
		Count     string `json:"count"`      // Кількість виконавчіх проваджень
		LiveCount string `json:"live_count"` // Кількість справ за якими заплановані засідання
	} `json:"administrative"`
	AdminOffense struct {
		Count     string `json:"count"`      // Кількість виконавчіх проваджень
		LiveCount string `json:"live_count"` // Кількість справ за якими заплановані засідання
	} `json:"admin_offense"`
}

// GetCompanyCourts
// Отримання кількості судових справ за видами судочинства, де компанія є стороною
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/companyCourtsCount
func (odb *OdbClient) GetCompanyCourts(
	code string, // код ЄДРПОУ компанії
) (response *CompanyCourtsList, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(companyCourtsEndpoint, map[string]string{
		"code": code,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "civil": {
	//    "count": "2",
	//    "live_count": "2"
	//  },
	//  "criminal": {
	//    "count": "2",
	//    "live_count": "2"
	//  },
	//  "arbitrage": {
	//    "count": "2",
	//    "live_count": "2"
	//  },
	//  "administrative": {
	//    "count": "2",
	//    "live_count": "2"
	//  },
	//  "admin_offense": {
	//    "count": "2",
	//    "live_count": "2"
	//  }
	//}
}

type CompanyCourtsDetail struct {
	Number           string `json:"number"`             // Номер
	Date             string `json:"date"`               // Датa
	DateStart        string `json:"date_start"`         // Датa
	LastScheduleDate string `json:"last_schedule_date"` // Дата останнього засідання
	Live             string `json:"live"`               // Ознака наявності засідань по справі в майбутньому
	Description      string `json:"description"`        // Суть справи
	ScheduleCount    string `json:"schedule_count"`     // Кількість засідань
	Cost             string `json:"cost"`               // Сума спору
	Amount           string `json:"amount"`             // Cума позовних вимог
	CourtName        string `json:"court_name"`         // Назва суду
	Plaintiffs       []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"plaintiffs"`
	Defendants []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"defendants"`
	ThirdPersons []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"third_persons"`
	Appeals []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"appeals"`
	Cassations []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"cassations"`
	JudgmentCode     string `json:"judgment_code"`      // Код типу судочинства
	LastDocumentDate string `json:"last_document_date"` // Дата останнього рішення
	Stages           struct {
		First struct {
			CourtCode            int    `json:"court_code"`             // Внутрішній код судової установи
			CourtName            string `json:"court_name"`             // Назва судової установи
			Judge                string `json:"judge"`                  // Суддя
			ConsiderationForSide string `json:"consideration_for_side"` // Результат
			Description          string `json:"description"`            // Опис результату рішення
		} `json:"first"`
		Appeal struct {
			CourtCode            int    `json:"court_code"`             // Внутрішній код судової установи
			CourtName            string `json:"court_name"`             // Назва судової установи
			Judge                string `json:"judge"`                  // Суддя
			ConsiderationForSide string `json:"consideration_for_side"` // Результат
			Description          string `json:"description"`            // Опис результату рішення
		} `json:"appeal"`
		Cassation struct {
			CourtCode            int    `json:"court_code"`             // Внутрішній код судової установи
			CourtName            string `json:"court_name"`             // Назва судової установи
			Judge                string `json:"judge"`                  // Суддя
			ConsiderationForSide string `json:"consideration_for_side"` // Результат
			Description          string `json:"description"`            // Опис результату рішення
		} `json:"cassation"`
	} `json:"stages"`
}

// GetCompanyCourtsByType
// Судові справи за типом судочинства, де компанія є стороною
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/companyCourtsDetail
func (odb *OdbClient) GetCompanyCourtsByType(
	courtsType string,
	code string,
	params map[string]string, //map[string]string{
	//	"sort_field":	"поле по якому відбувається сортування результату",
	//	"sort_type":	"порядок сортування (DESC - по зменшенню; ASC - по зростанню)",
	//	"date_from":	"фільтр з дати першого засідання або документа у справі",
	//	"date_to":		"фільтр по дату першого засідання або документа у справі",
	//	"offset":		"Зміщення відносно початку результатів пошуку",
	//	"limit":		"Кількість записів. Максимальний ліміт кількості записів — 1000",
	//	"date_from":	"Початкова дата пошуку (Y-m-d)",
	//	"date_to":		"Кінцева дата пошуку (Y-m-d)",
	//}
) (response *CompanyCourtsDetail, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(companyCourtsByTypeEndpoint, courtsType)

	params["code"] = code

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "number": "904/6017/17",
	//  "date": "2018-01-01",
	//  "date_start": "2018-01-01",
	//  "last_schedule_date": "2018-01-01",
	//  "live": "1",
	//  "description": "визнання недійсними рішень загальних зборів та договорів",
	//  "schedule_count": "1",
	//  "cost": "1",
	//  "amount": "1",
	//  "court_name": "Саксаганський районний суд м.Кривого Рогу",
	//  "plaintiffs": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "defendants": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "third_persons": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "appeals": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "cassations": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "judgment_code": "1",
	//  "last_document_date": "2018-01-01",
	//  "stages": {
	//    "first": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration_for_side": "lose",
	//      "description": "негативний для касанта"
	//    },
	//    "appeal": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration_for_side": "lose",
	//      "description": "негативний для касанта"
	//    },
	//    "cassation": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration_for_side": "lose",
	//      "description": "негативний для касанта"
	//    }
	//  }
	//}
}

type CompanyCourtsCases struct {
	Number           string `json:"number"`             // Номер
	Date             string `json:"date"`               // Дата
	DateStart        string `json:"date_start"`         // Дата
	LastScheduleDate string `json:"last_schedule_date"` // Дата останнього засідання
	LastStatus       string `json:"last_status"`        // Поточний стан розгляду справи
	Live             string `json:"live"`               // Ознака наявності засідань по справі в майбутньому
	Description      string `json:"description"`        // Суть справи
	ScheduleCount    string `json:"schedule_count"`     // Кількість засідань
	Cost             string `json:"cost"`               // Сума спору
	Amount           string `json:"amount"`             // Сума позовних вимог
	CourtName        string `json:"court_name"`         // Назва суду
	Plaintiffs       []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"plaintiffs"`
	Defendants []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"defendants"`
	ThirdPersons []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"third_persons"`
	Appeals []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"appeals"`
	Cassations []struct {
		Code string `json:"code"` // Код ЄДРПОУ
		Name string `json:"name"` // ПІБ
	} `json:"cassations"`
	JudgmentCode     string `json:"judgment_code"`      // Код типу судочинства
	LastDocumentDate string `json:"last_document_date"` // Дата останнього рішення
	Stages           struct {
		First struct {
			CourtCode     int    `json:"court_code"`    // Внутрішній код судової установи
			CourtName     string `json:"court_name"`    // Назва судової установи
			Judge         string `json:"judge"`         // Суддя
			Consideration string `json:"consideration"` // Результат
			Description   string `json:"description"`   // Опис результату рішення
			Decisions     []struct {
				CourtCode    int    `json:"court_code"`    // Внутрішній код судової установи
				CourtName    string `json:"court_name"`    // Назва судової установи
				JudgmentCode int    `json:"judgment_code"` // Внутрішній код Форми судочинства
				//Кримінальне
				//Цивільне
				//Господарське
				//Адміністративне
				//Адмінправопорушення
				JudgmentName string `json:"judgment_name"` // Форма судочинства
				JusticeCode  int    `json:"justice_code"`  // Внутрішній код Типу процесуального документа
				//Вирок
				//Постанова
				//Рішення
				//Судовий наказ
				//Ухвала
				//Окрема ухвала
				//Окрема думка
				JusticeName      string `json:"justice_name"`      // Тип процесуального документа
				AdjudicationDate string `json:"adjudication_date"` // Дата набрання законної сили
				DatePubl         string `json:"date_publ"`         // Дата публікації
				ReceiptDate      string `json:"receipt_date"`      // Дата реєстрації
				Judge            string `json:"judge"`             // Суддя
				Result           string `json:"result"`            // Результат
				Link             string `json:"link"`              // Посилання на рішення
			} `json:"decisions"`
		} `json:"first"`
		Appeal struct {
			CourtCode     int    `json:"court_code"`    // Внутрішній код судової установи
			CourtName     string `json:"court_name"`    // Назва судової установи
			Judge         string `json:"judge"`         // Суддя
			Consideration string `json:"consideration"` // Результат
			Description   string `json:"description"`   // Опис результату рішення
			Decisions     []struct {
				CourtCode    int    `json:"court_code"`    // Внутрішній код судової установи
				CourtName    string `json:"court_name"`    // Назва судової установи
				JudgmentCode int    `json:"judgment_code"` // Внутрішній код Форми судочинства
				//Кримінальне
				//Цивільне
				//Господарське
				//Адміністративне
				//Адмінправопорушення
				JudgmentName string `json:"judgment_name"` // Форма судочинства
				JusticeCode  int    `json:"justice_code"`  // Внутрішній код Типу процесуального документа
				//Вирок
				//Постанова
				//Рішення
				//Судовий наказ
				//Ухвала
				//Окрема ухвала
				//Окрема думка
				JusticeName      string `json:"justice_name"`      // Тип процесуального документа
				AdjudicationDate string `json:"adjudication_date"` // Дата набрання законної сили
				DatePubl         string `json:"date_publ"`         // Дата публікації
				ReceiptDate      string `json:"receipt_date"`      // Дата реєстрації
				Judge            string `json:"judge"`             // Суддя
				Result           string `json:"result"`            // Результат
				Link             string `json:"link"`              // Посилання на рішення
			} `json:"decisions"`
		} `json:"appeal"`
		Cassation struct {
			CourtCode     int    `json:"court_code"`    // Внутрішній код судової установи
			CourtName     string `json:"court_name"`    // Назва судової установи
			Judge         string `json:"judge"`         // Суддя
			Consideration string `json:"consideration"` // Результат
			Description   string `json:"description"`   // Опис результату рішення
			Decisions     []struct {
				CourtCode    int    `json:"court_code"`    // Внутрішній код судової установи
				CourtName    string `json:"court_name"`    // Назва судової установи
				JudgmentCode int    `json:"judgment_code"` // Внутрішній код Форми судочинства
				//Кримінальне
				//Цивільне
				//Господарське
				//Адміністративне
				//Адмінправопорушення
				JudgmentName string `json:"judgment_name"` // Форма судочинства
				JusticeCode  int    `json:"justice_code"`  // Внутрішній код Типу процесуального документа
				//Вирок
				//Постанова
				//Рішення
				//Судовий наказ
				//Ухвала
				//Окрема ухвала
				//Окрема думка
				JusticeName      string `json:"justice_name"`      // Тип процесуального документа
				AdjudicationDate string `json:"adjudication_date"` // Дата набрання законної сили
				DatePubl         string `json:"date_publ"`         // Дата публікації
				ReceiptDate      string `json:"receipt_date"`      // Дата реєстрації
				Judge            string `json:"judge"`             // Суддя
				Result           string `json:"result"`            // Результат
				Link             string `json:"link"`              // Посилання на рішення
			} `json:"decisions"`
		} `json:"cassation"`
	} `json:"stages"`
}

// GetCourtCases
// Отримання переліка інстанцій, рішень, позивачів, відповідачів, засідань за судовою справою. Отримання результату в кожній інстанції.
// https://docs.opendatabot.com/#/%D0%A1%D1%83%D0%B4%D0%BE%D0%B2%D0%B8%D0%B9%20%D1%80%D0%B5%D1%94%D1%81%D1%82%D1%80/CourtsCases
func (odb *OdbClient) GetCourtCases(
	number string,
	params map[string]string, //map[string]string{
	//	// 1 - Цивільні справи
	//	// 2 - Кримінальні справи
	//	// 3 - Господарські справи
	//	// 4 - Адміністративні справи
	//	// 5 - Справи про адмінправопорушення
	//	// Якщо параметр не зазначений,
	//	// а результат пошуку більше однієї справи з різним типом судочинства,
	//	// то виникне помилка неунікальності судової справи.
	//	// При зазначені типу судочинства, результат стає унікальним
	//	// та у відповіді відображається лише одна справа
	//	// Available values : 1, 2, 3, 4, 5
	//	"judgment_code": "Available values : 1, 2, 3, 4, 5",
	//}
) (response *CompanyCourtsCases, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(courtCasesEndpoint, number)

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "number": "904/6017/17",
	//  "date": "2018-01-01",
	//  "date_start": "2018-01-01",
	//  "last_schedule_date": "2018-01-01",
	//  "last_status": "Призначено до судового розгляду",
	//  "live": "1",
	//  "description": "визнання недійсними рішень загальних зборів та договорів",
	//  "schedule_count": "1",
	//  "cost": "1",
	//  "amount": "1",
	//  "court_name": "Саксаганський районний суд м.Кривого Рогу",
	//  "plaintiffs": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "defendants": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "third_persons": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "appeals": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "cassations": [
	//    {
	//      "code": "11111111",
	//      "name": "Петров Іван Володимирович"
	//    }
	//  ],
	//  "judgment_code": "1",
	//  "last_document_date": "2018-01-01",
	//  "stages": {
	//    "first": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration": "lose",
	//      "description": "негативний для касанта",
	//      "decisions": [
	//        {
	//          "court_code": 1521,
	//          "court_name": "Овідіопольський районний суд Одеської області",
	//          "judgment_code": 3,
	//          "judgment_name": "Кримінальне|Цивільне|Господарське|Адміністративне|Адмінправопорушення",
	//          "justice_code": 2,
	//          "justice_name": "Вирок|Постанова|Рішення|Судовий наказ|Ухвала|Окрема ухвала|Окрема думка",
	//          "adjudication_date": "2018-09-03 00:00:00+03",
	//          "date_publ": "2018-09-03 00:00:00+03",
	//          "receipt_date": "2018-09-03 00:00:00+03",
	//          "judge": "Попов І.В.",
	//          "result": "lose",
	//          "link": "https://opendatabot.com/court/76195512-b22b9fbfbb31d6aca70b89d1257287a4"
	//        }
	//      ]
	//    },
	//    "appeal": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration": "lose",
	//      "description": "негативний для касанта",
	//      "decisions": [
	//        {
	//          "court_code": 1521,
	//          "court_name": "Овідіопольський районний суд Одеської області",
	//          "judgment_code": 3,
	//          "judgment_name": "Кримінальне|Цивільне|Господарське|Адміністративне|Адмінправопорушення",
	//          "justice_code": 2,
	//          "justice_name": "Вирок|Постанова|Рішення|Судовий наказ|Ухвала|Окрема ухвала|Окрема думка",
	//          "adjudication_date": "2018-09-03 00:00:00+03",
	//          "date_publ": "2018-09-03 00:00:00+03",
	//          "receipt_date": "2018-09-03 00:00:00+03",
	//          "judge": "Попов І.В.",
	//          "result": "lose",
	//          "link": "https://opendatabot.com/court/76195512-b22b9fbfbb31d6aca70b89d1257287a4"
	//        }
	//      ]
	//    },
	//    "cassation": {
	//      "court_code": 1521,
	//      "court_name": "Овідіопольський районний суд Одеської області",
	//      "judge": "Шевченко Г. С.",
	//      "consideration": "lose",
	//      "description": "негативний для касанта",
	//      "decisions": [
	//        {
	//          "court_code": 1521,
	//          "court_name": "Овідіопольський районний суд Одеської області",
	//          "judgment_code": 3,
	//          "judgment_name": "Кримінальне|Цивільне|Господарське|Адміністративне|Адмінправопорушення",
	//          "justice_code": 2,
	//          "justice_name": "Вирок|Постанова|Рішення|Судовий наказ|Ухвала|Окрема ухвала|Окрема думка",
	//          "adjudication_date": "2018-09-03 00:00:00+03",
	//          "date_publ": "2018-09-03 00:00:00+03",
	//          "receipt_date": "2018-09-03 00:00:00+03",
	//          "judge": "Попов І.В.",
	//          "result": "lose",
	//          "link": "https://opendatabot.com/court/76195512-b22b9fbfbb31d6aca70b89d1257287a4"
	//        }
	//      ]
	//    }
	//  }
	//}
}

type Transports struct {
	Count int `json:"count"` // Кількість збігів
	Data  []struct {
		Id     int64  `json:"id"`     // Внутрішній id
		Number string `json:"number"` // Номер
	} `json:"data"`
}

// GetTransports
// Отримання переліку транспортних засобів
// https://docs.opendatabot.com/#/%D0%A2%D1%80%D0%B0%D0%BD%D1%81%D0%BF%D0%BE%D1%80%D1%82/transports
func (odb *OdbClient) GetTransports(
	params map[string]string, //map[string]string{
	//	"start":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//	"number":	"Номер транспортного засобу",
	//	"order":	"Порядок сортування (asc|desc)",
	//}
) (response *Transports, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(transportEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "count": 1,
	//  "data": [
	//    {
	//      "id": 3017223335761111,
	//      "number": "AA1122BB"
	//    }
	//  ]
	//}
}

type ItemFullTransport struct {
	Id            int    `json:"id"`           // Внутрішній id
	Number        string `json:"number"`       // Номер
	Model         string `json:"model"`        // Модель
	Year          string `json:"year"`         // Рік
	Date          string `json:"date"`         // Дата реєстрації
	Registration  string `json:"registration"` // Вид реєстрації
	Capacity      int    `json:"capacity"`     // Об'єм двигуна
	OwnerHash     string `json:"owner_hash"`   // Внутрішній id власника
	Color         string `json:"color"`        // Колір
	Kind          string `json:"kind"`         // Тип
	Body          string `json:"body"`         // Тип кузову
	OwnWeight     int    `json:"own_weight"`   // Вага
	RegAddrKoatuu string `json:"reg_addr_koatuu"`
	DepCode       string `json:"dep_code"`
	Dep           string `json:"dep"`
}

// GetTransportById
// Отримання інформації по реєстрації транспортного засобу
// https://docs.opendatabot.com/#/%D0%A2%D1%80%D0%B0%D0%BD%D1%81%D0%BF%D0%BE%D1%80%D1%82/transport
func (odb *OdbClient) GetTransportById(
	id string, // внутрішній id, який отримали при пошуку транспортних засобів
) (response *ItemFullTransport, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(transportByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "id": 1,
	//  "number": "AA1122BB",
	//  "model": "SSANG YONG REXTON",
	//  "year": "2007",
	//  "date": "2018-07-09",
	//  "registration": "40 - ВТОРИННА РЕЄСТРАЦІЯ ТЗ, ПРИДБАНОГО В ТОРГОВЕЛЬНІЙ ОРГАНІЗАЦІЇ",
	//  "capacity": 2345,
	//  "owner_hash": "asdnvevew21231red3r23Xqw",
	//  "color": "СІРИЙ",
	//  "kind": "ЛЕГКОВИЙ",
	//  "body": "СЕДАН-B",
	//  "own_weight": 1200,
	//  "reg_addr_koatuu": "8033400051",
	//  "dep_code": "1333408",
	//  "dep": "ВРЕР-8 УДАІ В М.КИЄВІ"
	//}
}

type TransportLicenses struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Id               int    `json:"id"`                 // Внутрішній id
			Number           string `json:"number"`             // Номер
			LicenseStatus    string `json:"license_status"`     // Статус ліцензії
			LicenseIssueDate string `json:"license_issue_date"` // дата випуску ліцензії
			LicenseStartDate string `json:"license_start_date"` // дата початку ліцензії
			LicenseEndDate   string `json:"license_end_date"`   // кінцева дата ліцензії
			LicenseType      string `json:"license_type"`       // Тимчасовий реєстраційний талон
		} `json:"items"`
	} `json:"data"`
}

// GetTransportLicenses
// Отримання переліку ліцензій транспортних засобів
// https://docs.opendatabot.com/#/%D0%A2%D1%80%D0%B0%D0%BD%D1%81%D0%BF%D0%BE%D1%80%D1%82/transportLicenses
func (odb *OdbClient) GetTransportLicenses(
	params map[string]string, //map[string]string{
	//	"offset":		"Зміщення відносно початку результатів пошуку",
	//	"limit":		"Кількість записів",
	//	"number":		"Номер транспортного засобу",
	//	"code":			"Код компанії або ІНН ФОП",
	//	"owner_hash":	"Внутрішній id власника",
	//}
) (response *TransportLicenses, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(transportLicensesEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "id": 1,
	//        "number": "BX7768XT",
	//        "license_status": "Роздрукована",
	//        "license_issue_date": "2016-02-23",
	//        "license_start_date": "2017-02-23",
	//        "license_end_date": "2022-02-23",
	//        "license_type": "Тимчасовий реєстраційний талон\t"
	//      }
	//    ]
	//  }
	//}
}

type ItemFullTransportLicenses struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Id               int    `json:"id"`                 // Внутрішній id
		Number           string `json:"number"`             // Номер
		CarrierName      string `json:"carrier_name"`       // Перевізник
		OwnerHash        string `json:"owner_hash"`         // Внутрішній id власника
		LicenseStatus    string `json:"license_status"`     // Статус ліцензії
		LicenseIssueDate string `json:"license_issue_date"` // дата випуску ліцензії
		LicenseStartDate string `json:"license_start_date"` // дата початку ліцензії
		LicenseEndDate   string `json:"license_end_date"`   // кінцева дата ліцензії
		LicenseType      string `json:"license_type"`       // Тимчасовий реєстраційний талон
		TransportType    string `json:"transport_type"`     // тип транспорту
		TransportStatus  string `json:"transport_status"`   // статус транспорту
		TransportVendor  string `json:"transport_vendor"`   // марка
		TransportModel   string `json:"transport_model"`    // модель
		TransportYear    string `json:"transport_year"`     // рік виробництва
		TransportSeats   string `json:"transport_seats"`    // посадочні місця
		Vin              string `json:"vin"`                // vin номер
		Code             string `json:"code"`               // Код компанії або ІПН ФОП
	} `json:"data"`
}

// GetTransportLicensesById
// Отримання інформації про ліцензію транспортного засобу
// https://docs.opendatabot.com/#/%D0%A2%D1%80%D0%B0%D0%BD%D1%81%D0%BF%D0%BE%D1%80%D1%82/transportLicense
func (odb *OdbClient) GetTransportLicensesById(
	id string, // внутрішній id, який отримали при пошуку ліцензій транспортних засобів
) (response *ItemFullTransportLicenses, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(transportLicensesByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "id": 1,
	//    "number": "BX7768XT",
	//    "carrier_name": "ПП ТРАНС-АВТО-Д",
	//    "owner_hash": "f784c979c15e551defc1727257ad0351",
	//    "license_status": "Роздрукована",
	//    "license_issue_date": "2016-02-23",
	//    "license_start_date": "2017-02-23",
	//    "license_end_date": "2022-02-23",
	//    "license_type": "Тимчасовий реєстраційний талон\t",
	//    "transport_type": "Автобус",
	//    "transport_status": "Звичайний",
	//    "transport_vendor": "MERCEDES-BENZ",
	//    "transport_model": "312",
	//    "transport_year": "1997",
	//    "transport_seats": "18",
	//    "vin": "WDB9055631P239719",
	//    "code": "00214244"
	//  }
	//}
}

type GenKey struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		ApiKey        string `json:"apiKey"`         // Згенерований API ключ
		SettingsToken string `json:"settings_token"` // токен
	} `json:"data"`
}

// GetGenKey
// Генерація API ключа для кінцевого користувача партнера
// https://docs.opendatabot.com/#/%D0%A0%D0%BE%D0%B1%D0%BE%D1%82%D0%B0%20API/apiKeyGeneration
func (odb *OdbClient) GetGenKey(
	salt string, // пароль партнера
	id string, // незмінний внутрішній ідентифікатор клієнта, строка або число
) (response *GenKey, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(genKeyEndpoint, map[string]string{
		"salt": salt,
		"id":   id,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "apiKey": "ghpCpWUS-2345353-603092018de9c199ca8a0209efe40f27",
	//    "settings_token": "f1f0e8750f21627cbae87918d203400e"
	//  }
	//}
}

type Statistics struct {
	COMPANY struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"COMPANY"`
	FULLCOMPANY struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"FULLCOMPANY"`
	FOP struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"FOP"`
	FOPINN struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"FOPINN"`
	PERSON struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"PERSON"`
	REGISTRATIONS struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"REGISTRATIONS"`
	VAT struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"VAT"`
	SCHEDULE struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"SCHEDULE"`
	COMPANYRECORD struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"COMPANYRECORD"`
	COURT struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"COURT"`
	SUBSCRIPTION struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"SUBSCRIPTION"`
	UNSUBSCRIPTION struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"UNSUBSCRIPTION"`
	HISTORY struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"HISTORY"`
	CHANGES struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"CHANGES"`
	INSTITUTIONS struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"INSTITUTIONS"`
	SEARCH struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"SEARCH"`
	LISTS struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"LISTS"`
	DEBT struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"DEBT"`
	APICOURT struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"APICOURT"`
	MESSAGE struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"MESSAGE"`
	STATISTICS struct {
		Name    string `json:"name"`    // Назва
		Used    int    `json:"used"`    // Використано запитів
		Limit   int    `json:"limit"`   // Кількість записів
		Balance int    `json:"balance"` // Поточний баланс запитів
	} `json:"STATISTICS"`
	ExpiryDate string `json:"expiry_date"` // Дата закінчення пакету
	CustomerId string `json:"customerId"`  // ID клієнта
	Webhook    string `json:"webhook"`     // Встановленний webhook
}

// GetStatistics
// Отримання повної інформації про використання API запитів
// https://docs.opendatabot.com/#/%D0%A0%D0%BE%D0%B1%D0%BE%D1%82%D0%B0%20API/statistics
func (odb *OdbClient) GetStatistics() (response *Statistics, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(statisticsEndpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "COMPANY": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "FULLCOMPANY": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "FOP": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "FOPINN": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "PERSON": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "REGISTRATIONS": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "VAT": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "SCHEDULE": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "COMPANYRECORD": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "COURT": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "SUBSCRIPTION": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "UNSUBSCRIPTION": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "HISTORY": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "CHANGES": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "INSTITUTIONS": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "SEARCH": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "LISTS": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "DEBT": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "APICOURT": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "MESSAGE": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "STATISTICS": {
	//    "name": "Базове API",
	//    "used": 56,
	//    "limit": 100,
	//    "balance": 44
	//  },
	//  "expiry_date": "2019-10-10",
	//  "customerId": "168",
	//  "webhook": "https://opendatabot.ua/webhook"
	//}
}

type AlimentData struct {
	Count    int `json:"count"` // Кількість збігів
	Aliments []struct {
		FullName  string `json:"full_name"`  // Повне ім'я
		BirthDate string `json:"birth_date"` // Дата народження
		Active    int    `json:"active"`     // Ознака актуальності
	} `json:"aliments"`
}

// GetAliment
// Отримання публічної інформації щодо особи, наявність в базі боржників за аліментами
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/aliment
func (odb *OdbClient) GetAliment(
	pib string,
	params map[string]string, //map[string]string{
	//	"start":		"Зміщення відносно початку результатів пошуку",
	//	"birth_date":	"Фільтр за датою народження в форматі",
	//	"limit":		"Кількість записів",
	//}
) (response *AlimentData, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	params["pib"] = pib

	err = odb.Do(alimentEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "count": 1,
	//  "aliments": [
	//    {
	//      "full_name": "Шевченко Олександр Володимирович",
	//      "birth_date": "1970-01-01",
	//      "active": 1
	//    }
	//  ]
	//}
}

type Lawyers struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Id           int    `json:"id"`            // Внутрішній id
			FullName     string `json:"full_name"`     // ПІБ
			Racalc       string `json:"racalc"`        // Обліковується у
			Certnum      string `json:"certnum"`       // № Свідоцтва
			Certat       string `json:"certat"`        // Дата видачі свідоцтва
			Certcalc     string `json:"certcalc"`      // Орган, що видав свідоцтво
			DatabaseDate string `json:"database_date"` // Дата актуальності
		} `json:"items"`
	} `json:"data"`
}

// GetLawyers
// Отримання переліку адвокатів
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/lawyers
func (odb *OdbClient) GetLawyers(
	params map[string]string, //map[string]string{
	//	"offset":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//	"name":		"ПІБ особи",
	//}
) (response *Lawyers, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(lawyersEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "id": 1,
	//        "full_name": "Петров Іван Володимирович",
	//        "racalc": "Рада адвокатів Дніпропетровської області",
	//        "certnum": "3462",
	//        "certat": "2018-11-13 00:00:00",
	//        "certcalc": "Рада адвокатів Дніпропетровської області",
	//        "database_date": "2018-11-15"
	//      }
	//    ]
	//  }
	//}
}

type Lawyer struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Id             int    `json:"id"`              // Внутрішній id
		FullName       string `json:"full_name"`       // ПІБ
		Racalc         string `json:"racalc"`          // Обліковується у
		Certnum        string `json:"certnum"`         // № Свідоцтва
		Certat         string `json:"certat"`          // Дата видачі свідоцтва
		Certcalc       string `json:"certcalc"`        // Орган, що видав свідоцтво
		DatabaseDate   string `json:"database_date"`   // Дата актуальності
		Phone          string `json:"phone"`           // Мобільний
		Email          string `json:"email"`           // E-mail
		DecisionDate   string `json:"decision_date"`   // Дата прийняття рішення
		DecisionNumber string `json:"decision_number"` // Номер рішення
		Activities     string `json:"activities"`      // Форми адвокатської діяльності
		Experience     string `json:"experience"`      // Загальний стаж адвоката
		Termination    string `json:"termination"`     // Інформація про зупинення або припинення права на заняття адвокатською діяльністю
	} `json:"data"`
}

// GetLawyerById
// Отримання інформації про Адвоката
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/lawyer
func (odb *OdbClient) GetLawyerById(
	id string, // внутрішній id, який отримали при пошуку адвокатів
) (response *Lawyer, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(lawyersByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "id": 1,
	//    "full_name": "Петров Іван Володимирович",
	//    "racalc": "Рада адвокатів Дніпропетровської області",
	//    "certnum": "3462",
	//    "certat": "2018-11-13 00:00:00",
	//    "certcalc": "Рада адвокатів Дніпропетровської області",
	//    "database_date": "2018-11-15",
	//    "phone": "+38(093)459-65-14",
	//    "email": "email@gmail.com",
	//    "decision_date": "2018-11-09",
	//    "decision_number": "924",
	//    "activities": "Індивідуальна адвокатська діяльність",
	//    "experience": "з 1984 року",
	//    "termination": "Право на заняття адвокатською діяльністю зупинено згідно п.1 ч.1 ст.31 ЗУ 'Про адвокатуру та адвокатську діяльність' з 09.11.2018 на підставі заяви адвоката"
	//  }
	//}
}

type CorruptOfficialsItem struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Id             string   `json:"id"`              // ID
		FullName       string   `json:"full_name"`       // Повне ім'я
		DecisionDate   string   `json:"decision_date"`   // Дата судового рішення
		DecisionNumber string   `json:"decision_number"` // Номер судового рішення
		WorkPlace      string   `json:"work_place"`      // Місце роботи на час вчинення корупційного правопорушення
		Position       string   `json:"position"`        // Посада на час вчинення корупційного правопорушення
		CodexArticles  []string `json:"codex_articles"`  // Статті кодексів
		Active         int      `json:"active"`          // Ознака актуальності
	} `json:"data"`
}

// GetCorruptOfficialsById
// Отримання відомостей про осібу, яка вчинила корупційні правопорушення, за внутрішнім id
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/corrupt-officials-item
func (odb *OdbClient) GetCorruptOfficialsById(
	id string, // внутрішній id, який отримали при пошуку корупціонерів по ПІБ
) (response *CorruptOfficialsItem, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(corruptOfficialsByIdEndpoint, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "id": "d20459dba9295d1e18abb0a7065198ec",
	//    "full_name": "Шевченко Олександр Володимирович",
	//    "decision_date": "2011-01-01",
	//    "decision_number": "68483631",
	//    "work_place": "Трапівська сільська рада",
	//    "position": "депутат",
	//    "codex_articles": [
	//      "Кодекс України про адміністративні правопорушення  Стаття 172-7. Порушення вимог щодо повідомлення про конфлікт інтересів"
	//    ],
	//    "active": 1
	//  }
	//}
}

type CorruptOfficials struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Id             string   `json:"id"`              // ID
			FullName       string   `json:"full_name"`       // Повне ім'я
			DecisionDate   string   `json:"decision_date"`   // Дата судового рішення
			DecisionNumber string   `json:"decision_number"` // Номер судового рішення
			WorkPlace      string   `json:"work_place"`      // Місце роботи на час вчинення корупційного правопорушення
			Position       string   `json:"position"`        // Посада на час вчинення корупційного правопорушення
			CodexArticles  []string `json:"codex_articles"`  // Статті кодексів
			Active         int      `json:"active"`          // Ознака актуальності
		} `json:"items"`
	} `json:"data"`
}

// GetCorruptOfficials
// Отримання відомостей про осіб, які вчинили корупційні правопорушення
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/corrupt-officials
func (odb *OdbClient) GetCorruptOfficials(
	pib string, // ПІБ особи
	params map[string]string, // map[string]string{
	//	"start":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//}
) (response *CorruptOfficials, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	params["pib"] = pib

	err = odb.Do(corruptOfficialsEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "id": "d20459dba9295d1e18abb0a7065198ec",
	//        "full_name": "Шевченко Олександр Володимирович",
	//        "decision_date": "2011-01-01",
	//        "decision_number": "68483631",
	//        "work_place": "Трапівська сільська рада",
	//        "position": "депутат",
	//        "codex_articles": [
	//          "Кодекс України про адміністративні правопорушення  Стаття 172-7. Порушення вимог щодо повідомлення про конфлікт інтересів"
	//        ],
	//        "active": 1
	//      }
	//    ]
	//  }
	//}
}

type Passport struct {
	Count int `json:"count"` // Кількість збігів
	Data  []struct {
		Id        string `json:"id"`         // ID запису в МВС України
		Number    string `json:"number"`     // Номер паспорта
		Type      string `json:"type"`       // invalid, lost
		Ovd       string `json:"ovd"`        // Регіон (орган внутрішніх справ)
		TheftDate string `json:"theft_date"` // Дата внесення в базу
		Date      string `json:"date"`       // Дата внесення змін
	} `json:"data"`
}

// GetPassport
// Отримання інформації про викрадені/втрачені паспорти громадянина України
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/passport
func (odb *OdbClient) GetPassport(
	number string, // Номер паспорту, наприклад CP634742
) (response *Passport, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(passportEndpoint, map[string]string{
		"number": number,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "count": 1,
	//  "data": [
	//    {
	//      "id": "3015178477784375",
	//      "number": "ME394233",
	//      "type": "lost",
	//      "ovd": "ОБОЛОНСЬКЕ УПРАВЛІННЯ ПОЛІЦІЇ ГУНП В М. КИЄВІ",
	//      "theft_date": "2015-06-27 13:15:07",
	//      "date": "2015-06-27 13:15:07"
	//    }
	//  ]
	//}
}

type Wanted struct {
	Status string `json:"status"` // Кількість збігів
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Id          string `json:"id"`           // Внутрішній ідентифікатор МВС
			FullName    string `json:"full_name"`    // ім'я
			BirthDate   string `json:"birth_date"`   // дата народження
			LostDate    string `json:"lost_date"`    // Дата пошуку
			Sex         string `json:"sex"`          // Стать
			ArticleCrim string `json:"article_crim"` // звинувачення
			LostPlace   string `json:"lost_place"`
			Ovd         string `json:"ovd"` // розшукує
			Category    string `json:"category"`
			Restraint   string `json:"restraint"`   // Запобіжний захід
			StatusText  string `json:"status_text"` // текст
			Status      string `json:"status"`      // статус
		} `json:"items"`
	} `json:"data"`
}

// GetWanted
// Отримання інформації по базі людей в розшуку
// https://docs.opendatabot.com/#/%D0%A4%D1%96%D0%B7%D0%B8%D1%87%D0%BD%D1%96%20%D0%BE%D1%81%D0%BE%D0%B1%D0%B8/wanted
func (odb *OdbClient) GetWanted(
	pib string, // ПІБ особи
	params map[string]string, // map[string]string{
	//	"start":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//}
) (response *Wanted, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	params["pib"] = pib

	err = odb.Do(wantedEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "id": "3023451632343695",
	//        "full_name": "Шевченко Олександр Ігорович",
	//        "birth_date": "1988-03-08",
	//        "lost_date": "2019-06-22",
	//        "sex": "male",
	//        "article_crim": "СТ.128 Ч.1",
	//        "lost_place": "ОДЕССКАЯ, КИЕВСКИЙ, ОДЕССА",
	//        "ovd": "ТАЇРОВСЬКЕ ВІДДІЛЕННЯ ПОЛІЦІЇ КИЇВСЬКОГО ВІДДІЛУ ПОЛІЦІЇ В М. ОДЕСІ ГУНП В ОДЕСЬКІЙ ОБЛАСТІ",
	//        "category": "особа, яка переховується від органів досудового розслідування",
	//        "restraint": "ухвала суду про дозвіл на затримання  з метою приводу",
	//        "status_text": "В розшуку з 23.06.2017",
	//        "status": "статус"
	//      }
	//    ]
	//  }
	//}
}

type FullPenaltiesSuccess struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count       int `json:"count"`        // Кількість збігів
		ActiveCount int `json:"active_count"` // Кількість збігів
		Items       []struct {
			Number             string `json:"number"`               // Номер виконавчого провадження
			BorrowerCode       string `json:"borrower_code"`        // код ЄДРПОУ
			SubType            string `json:"sub_type"`             // Тип боржника
			BorrowerLastName   string `json:"borrower_last_name"`   // Прізвище боржника
			BorrowerFirstName  string `json:"borrower_first_name"`  // Ім'я боржника
			BorrowerMiddleName string `json:"borrower_middle_name"` // По-батькові боржника
			BorrowerBirthDate  string `json:"borrower_birth_date"`  // Дата народження боржника
			CreditorName       string `json:"creditor_name"`        // Найменування стягувача
			CreditorCode       string `json:"creditor_code"`        // Код ЄДРПОУ стягувача
			CreditorSubType    string `json:"creditor_sub_type"`    // Тип стягувача
			AsvpGisName        string `json:"asvp_gis_name"`        // Орган ДВС
			AsvpDepId          string `json:"asvp_dep_id"`          // Ідентіфікаційний номер Органа ДВС
			BeginDate          string `json:"begin_date"`           // Дата відкриття провадження
			AsvpStatus         string `json:"asvp_status"`          // Статус провадження
			Active             string `json:"active"`               // Код статусу провадження
		} `json:"items"`
	} `json:"data"`
}

// GetFullPenaltyByNumber
// Отримання інформації про історію виконавчих проваджень компанії або приватної особи за номером провадження
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/full_penalties_number
func (odb *OdbClient) GetFullPenaltyByNumber(
	number string, // Номер виконавчого провадження
	params map[string]string, // map[string]string{
	//	"source":	"Джерело з якого виконується витяг інформації по виконавчим провадженням, opendatabot - для отримання інформації з бази даних Opendatabot",
	//}
) (response *FullPenaltiesSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(fullPenaltyByNumberEndpoint, number)

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "active_count": 1,
	//    "items": [
	//      {
	//        "number": "58526890",
	//        "borrower_code": "33433293",
	//        "sub_type": "Фізична особа",
	//        "borrower_last_name": "Агаркова",
	//        "borrower_first_name": "Діана",
	//        "borrower_middle_name": "Володимирівна",
	//        "borrower_birth_date": "1989-05-13",
	//        "creditor_name": "Міське комунальне підприємство «Миколаївводоканал»",
	//        "creditor_code": "31448144",
	//        "creditor_sub_type": "Юридична особа",
	//        "asvp_gis_name": "Приватний виконавець Куліченко Д.О.",
	//        "asvp_dep_id": "80718",
	//        "begin_date": "2019-03-04 00:00:00",
	//        "asvp_status": "Примусове виконання",
	//        "active": "1"
	//      }
	//    ]
	//  }
	//}
}

type FullPenaltiesSecretSuccess struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Number             string `json:"number"`               // Номер виконавчого провадження
		BorrowerCode       string `json:"borrower_code"`        // код ЄДРПОУ
		SubType            string `json:"sub_type"`             // Тип боржника
		BorrowerLastName   string `json:"borrower_last_name"`   // Прізвище боржника
		BorrowerFirstName  string `json:"borrower_first_name"`  // Ім'я боржника
		BorrowerMiddleName string `json:"borrower_middle_name"` // По-батькові боржника
		BorrowerBirthDate  string `json:"borrower_birth_date"`  // Дата народження боржника
		CreditorName       string `json:"creditor_name"`        // Найменування стягувача
		CreditorCode       string `json:"creditor_code"`        // Код ЄДРПОУ стягувача
		CreditorSubType    string `json:"creditor_sub_type"`    // Тип стягувача
		AsvpGisName        string `json:"asvp_gis_name"`        // Орган ДВС
		AsvpDepId          string `json:"asvp_dep_id"`          // Ідентіфікаційний номер Органа ДВС
		BeginDate          string `json:"begin_date"`           // Дата відкриття провадження
		AsvpStatus         string `json:"asvp_status"`          // Статус провадження
		Active             string `json:"active"`               // Код статусу провадження
		State              string `json:"state"`                // Статус
		ExecutorName       string `json:"executor_name"`        // П.І.Б виконавця
		Publisher          string `json:"publisher"`            // Орган, який видав виконавчий документ
		PublisherInfo      string `json:"publisher_info"`       // Дата та номер виконавчого документу
		ExecutorAdress     string `json:"executor_adress"`      // Адреса виконавця
		Documents          []struct {
			Id         string `json:"id"`          // Ідентифікаційний номер документа
			Name       string `json:"name"`        // Назва документа
			PrintDate  string `json:"print_date"`  // Дата публікації
			AcceptDate string `json:"accept_date"` // Дата прийняття
			CancelDate string `json:"cancel_date"` // Дата скасування документу
			Link       string `json:"link"`        // Посилання на документ
		} `json:"documents"`
	} `json:"data"`
}

// GetFullPenaltyDocByNumber
// Отримання інформації про історію виконавчих проваджень компанії або приватної особи за номером провадження та ідентифікатором доступу
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/full_penalties_secret
func (odb *OdbClient) GetFullPenaltyDocByNumber(
	number string, // Номер виконавчого провадження
	secret string, // Ідентифікатор доступу
) (response *FullPenaltiesSecretSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(fullPenaltyDocByNumberEndpoint, number)

	params := map[string]string{
		"secret": secret,
	}

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "number": "58526890",
	//    "borrower_code": "33433293",
	//    "sub_type": "Фізична особа",
	//    "borrower_last_name": "Агаркова",
	//    "borrower_first_name": "Діана",
	//    "borrower_middle_name": "Володимирівна",
	//    "borrower_birth_date": "1989-05-13",
	//    "creditor_name": "Міське комунальне підприємство «Миколаївводоканал»",
	//    "creditor_code": "31448144",
	//    "creditor_sub_type": "Юридична особа",
	//    "asvp_gis_name": "Приватний виконавець Куліченко Д.О.",
	//    "asvp_dep_id": "80718",
	//    "begin_date": "2019-03-04 00:00:00",
	//    "asvp_status": "Примусове виконання",
	//    "active": "1",
	//    "state": "Прийнятий до виконання",
	//    "executor_name": "Шевченка Павло Іванович",
	//    "publisher": "Свалявський районний суд",
	//    "publisher_info": "виконавчий лист 27.09.2016 №6/306/48/16",
	//    "executor_adress": "Закарпатська обл.",
	//    "documents": [
	//      {
	//        "id": "12086725532",
	//        "name": "Постанова про відкриття виконавчого  провадження (з ідентифікатором)",
	//        "print_date": "2018-09-10 10:58:59",
	//        "accept_date": "2018-09-10 10:58:59",
	//        "cancel_date": "2018-09-10",
	//        "link": "Закарпатська обл."
	//      }
	//    ]
	//  }
	//}
}

// GetFullPenalty
// Отримання інформації про історію виконавчих проваджень компанії або приватної особи за стороною провадження
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/full_penalties_params
func (odb *OdbClient) GetFullPenalty(
	params map[string]string, // map[string]string{
	//	"borrower_code":		"код ЄДРПОУ боржника",
	//	"creditor_code":		"код ЄДРПОУ стягувача",
	//	"borrower_first_name":	"Ім'я боржника",
	//	"borrower_last_name":	"Прізвище боржника",
	//	"borrower_middle_name":	"По-батькові боржника",
	//	"borrower_birth_date":	"Дата народження боржника",
	//	"offset":				"Зміщення відносно початку результатів пошуку",
	//	"limit":				"Кількість записів",
	//	"source":				"Джерело з якого виконується витяг інформації по виконавчим провадженням, opendatabot - для отримання інформації з бази даних Opendatabot",
	//}
) (response *FullPenaltiesSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(fullPenaltyEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "active_count": 1,
	//    "items": [
	//      {
	//        "number": "58526890",
	//        "borrower_code": "33433293",
	//        "sub_type": "Фізична особа",
	//        "borrower_last_name": "Агаркова",
	//        "borrower_first_name": "Діана",
	//        "borrower_middle_name": "Володимирівна",
	//        "borrower_birth_date": "1989-05-13",
	//        "creditor_name": "Міське комунальне підприємство «Миколаївводоканал»",
	//        "creditor_code": "31448144",
	//        "creditor_sub_type": "Юридична особа",
	//        "asvp_gis_name": "Приватний виконавець Куліченко Д.О.",
	//        "asvp_dep_id": "80718",
	//        "begin_date": "2019-03-04 00:00:00",
	//        "asvp_status": "Примусове виконання",
	//        "active": "1"
	//      }
	//    ]
	//  }
	//}
}

type PerformerSuccess struct {
	Status string `json:"status"`
	Data   struct {
		Count string `json:"count"`
		Items []struct {
			RegionId string `json:"regionId"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			Address  string `json:"address"`
			Contacts string `json:"contacts"`
			Managers string `json:"managers"`
		} `json:"items"`
	} `json:"data"`
}

// GetPerformer
// Отримання інформації про державні та приватні виконавчі служби
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/performer
func (odb *OdbClient) GetPerformer(
	params map[string]string, // map[string]string{
	//	"name":			"Назва або ПІБ виконавчої служби",
	//	"region_id":	"Ідентифікатор регіону:", //1 - Автономна Республіка Крим
	//	//2 - Вінницька обл
	//	//3 - Волинська обл
	//	//4 - Дніпропетровська обл
	//	//5 - Донецька обл
	//	//6 - Житомирська обл
	//	//7 - Закарпатська обл
	//	//8 - Запорізька обл
	//	//9 - Івано-Франківська обл
	//	//10 - Київська обл
	//	//11 - Кіровоградська обл
	//	//12 - Луганська обл
	//	//13 - Львівська обл
	//	//14 - Миколаївська обл
	//	//15 - Одеська обл
	//	//16 - Полтавська обл
	//	//17 - Рівненська обл
	//	//18 - Сумська обл
	//	//19 - Тернопільська обл
	//	//20 - Харківська обл
	//	//21 - Херсонська обл
	//	//22 - Хмельницька обл
	//	//23 - Черкаська обл
	//	//24 - Чернівецька обл
	//	//25 - Чернігівська обл
	//	//26 - м.Київ
	//	//27 - м.Севастополь
	//	"type":		"Державна або приватна виконавча служба. Available values: private, government",
	//	"offset":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//}
) (response *PerformerSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(performerEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "items": [
	//      {
	//        "regionId": "2",
	//        "name": "Ямпільський районний відділ державної виконавчої служби Головного територіального управління юстиції у Вінницькій області",
	//        "type": "private",
	//        "address": "Старокозацька, 56, м. Дніпро, 49101",
	//        "contacts": "info_prim@dp.dvs.gov.ua (056) 778-07-45, 778-07-46, 778-09-22, 778-09-20 (056) 778-09-22",
	//        "managers": "Півень Сергій Володимирович (056) 778-09-22 info_prim@dp.dvs.gov.ua Міхіна Ольга Іванівна (056) 778-09-20 info_prim@dp.dvs.gov.ua"
	//      }
	//    ]
	//  }
	//}
}

type PenaltiesSuccess struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"` // Кількість збігів
		Items []struct {
			Code             string    `json:"code"`                // код ЄДРПОУ
			CourtName        string    `json:"court_name"`          // Документ виданий
			GisName          string    `json:"gis_name"`            // Зв'язок з виконавцем
			Number           string    `json:"number"`              // Номер виконавчого провадження
			Category         string    `json:"category"`            // Категорія стягнення
			Id               string    `json:"id"`                  // Ідентіфікаційний номер
			Name             string    `json:"name"`                // Назва
			AddressAtuStr    string    `json:"address_atu_str"`     // Адреса виконавця
			Address          string    `json:"address"`             // Адреса виконавця
			DepartmentPhone  string    `json:"department_phone"`    // Номер телефону виконавця
			Executor         string    `json:"executor"`            // Виконавець
			ExecutorPhone    string    `json:"executor_phone"`      // Номер телефону виконавця
			ExecutorEmail    string    `json:"executor_email"`      // Email виконавця
			DeductionType    string    `json:"deduction_type"`      // Категорія стягнення
			LastName         string    `json:"last_name"`           // Прізвище боржника
			FirstName        string    `json:"first_name"`          // Ім'я боржника
			MiddleName       string    `json:"middle_name"`         // Ім'я по батькові боржника
			BirthDate        time.Time `json:"birth_date"`          // Дата народження боржника
			BirthPlaceAtuStr string    `json:"birth_place_atu_str"` // Місце народження боржника
			BirthPlace       string    `json:"birth_place"`         // Адреса народження боржника
			Link             string    `json:"link"`                // Посилання на додаткову інформацію
		} `json:"items"`
	} `json:"data"`
}

// GetPenaltiesByCode
// Отримання інформації про актуальні виконавчі провадження компанії або приватної особи за кодом боржника
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/penalties
func (odb *OdbClient) GetPenaltiesByCode(
	code string, // код ЄДРПОУ
	params map[string]string, // map[string]string{
	//	"categories[1]":	"Код категорії", //01 - стягнення коштів
	//  //02 - звернення стягнення на майно
	//  //03 - стягнення аліментів
	//  //04 - стягнення періодичних платежів (крім аліментів)
	//  //05 - стягнення заборгованості із заробітної плати та інших платежів, пов’язаних з трудовими відносинами
	//  //06 - стягнення соціальних виплат
	//  //07 - стягнення заборгованості з оплати комунальних послуг
	//  //08 - стягнення штрафів у справах про адміністративні правопорушення
	//  //09 - стягнення штрафів у справах про адміністративні правопорушення у сфері безпеки дорожнього руху
	//  //10 - забезпечення позову
	//  //11 - зобов’язання вчинити певні дії або утриматися від їх вчинення
	//  //12 - поновлення на роботі
	//  //13 - вселення стягувача
	//  //14 - виселення
	//  //15 - відібрання дитини
	//  //16 - заборона вчиняти певні дії
	//  //17 - конфіскація майна
	//  //18 - конфіскація майна, вилученого митними органами
	//  //18.1 - конфіскація майна засуджених
	//  //19 - конфіскація коштів та майна за вчинення корупційного та пов’язаного з корупцією правопорушення
	//  //20 - оплатне вилучення
	//  //21 - передача стягувачу предметів, зазначених у виконавчому документі
	//  //22 - стягнення коштів на користь держави
	//  //23 - рішення Європейського суду з прав людини
	//  //24 - стягнення виконавчого збору
	//  //25 - стягнення витрат виконавчого провадження
	//  //26 - стягнення штрафів, накладених державним, приватним виконавцем
	//  //27 - стягнення основної винагороди приватного виконавця
	//  //28 - усунення перешкод у побаченні з дитиною, встановлення побачення з дитиною
	//	"offset":			"Зміщення відносно початку результатів пошуку",
	//	"limit":			"Кількість записів",
	//}
) (response *PenaltiesSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(penaltiesByCodeEndpoint, code)

	err = odb.Do(endpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "code": "33433293",
	//        "court_name": "суд загальної юрисдикції Селидівський міський суд Донецької області.Суддя:І.М.Владимирська",
	//        "gis_name": "Селидівський міський відділ державної виконавчої служби Головного територіального управління юстиції у Донецькій області",
	//        "number": "57654570",
	//        "category": "стягнення коштів на користь держави",
	//        "id": "12570284040",
	//        "name": "ДП Селидіввугілля",
	//        "address_atu_str": "Ворошилова 21",
	//        "address": "Ворошилова 21",
	//        "department_phone": "(06237) 9-08-09",
	//        "executor": "Глухівський міськрайонний відділ державної виконавчої служби",
	//        "executor_phone": "(06237) 9-08-09",
	//        "executor_email": "example@gmail.com",
	//        "deduction_type": "стягнення аліментів",
	//        "last_name": "Ткаченко",
	//        "first_name": "Олег",
	//        "middle_name": "Васильович",
	//        "birth_date": "1975-05-06T00:00:00Z",
	//        "birth_place_atu_str": "Дниіпро",
	//        "birth_place": "Короленка 10",
	//        "link": "https://opendatabot.com/api/v2/penalty/57654570?apiKey=<apiKey>"
	//      }
	//    ]
	//  }
	//}
}

type PenaltySuccess struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Code             string    `json:"code"`                // код ЄДРПОУ
		CourtName        string    `json:"court_name"`          // Документ виданий
		GisName          string    `json:"gis_name"`            // Зв'язок з виконавцем
		Number           string    `json:"number"`              // Номер виконавчого провадження
		Category         string    `json:"category"`            // Категорія стягнення
		Id               string    `json:"id"`                  // Ідентіфікаційний номер
		Name             string    `json:"name"`                // Назва
		AddressAtuStr    string    `json:"address_atu_str"`     // Адреса виконавця
		Address          string    `json:"address"`             // Адреса виконавця
		DepartmentPhone  string    `json:"department_phone"`    // Номер телефону виконавця
		Executor         string    `json:"executor"`            // Виконавець
		ExecutorPhone    string    `json:"executor_phone"`      // Номер телефону виконавця
		ExecutorEmail    string    `json:"executor_email"`      // Email виконавця
		DeductionType    string    `json:"deduction_type"`      // Категорія стягнення
		LastName         string    `json:"last_name"`           // Прізвище боржника
		FirstName        string    `json:"first_name"`          // Ім'я боржника
		MiddleName       string    `json:"middle_name"`         // Ім'я по батькові боржника
		BirthDate        time.Time `json:"birth_date"`          // Дата народження боржника
		BirthPlaceAtuStr string    `json:"birth_place_atu_str"` // Місце народження боржника
		BirthPlace       string    `json:"birth_place"`         // Адреса народження боржника
	} `json:"data"`
}

// GetPenaltyByNumber
// Отримання інформації про актуальні виконавчі провадження компанії або приватної особи за номером провадження
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/penalty
func (odb *OdbClient) GetPenaltyByNumber(
	number string, // Виконавчий номер
) (response *PenaltySuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(penaltyByNumberEndpoint, number)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "code": "33433293",
	//    "court_name": "суд загальної юрисдикції Селидівський міський суд Донецької області.Суддя:І.М.Владимирська",
	//    "gis_name": "Селидівський міський відділ державної виконавчої служби Головного територіального управління юстиції у Донецькій області",
	//    "number": "57654570",
	//    "category": "стягнення коштів на користь держави",
	//    "id": "12570284040",
	//    "name": "ДП Селидіввугілля",
	//    "address_atu_str": "Ворошилова 21",
	//    "address": "Ворошилова 21",
	//    "department_phone": "(06237) 9-08-09",
	//    "executor": "Глухівський міськрайонний відділ державної виконавчої служби",
	//    "executor_phone": "(06237) 9-08-09",
	//    "executor_email": "example@gmail.com",
	//    "deduction_type": "стягнення аліментів",
	//    "last_name": "Ткаченко",
	//    "first_name": "Олег",
	//    "middle_name": "Васильович",
	//    "birth_date": "1975-05-06T00:00:00Z",
	//    "birth_place_atu_str": "Дниіпро",
	//    "birth_place": "Короленка 10"
	//  }
	//}
}

type PenaltyByFioSuccess struct {
	Status string `json:"status"` // Статус операції
	Data   struct {
		Count int `json:"count"`
		Items []struct {
			CourtName       string    `json:"court_name"`       // Документ виданий
			GisName         string    `json:"gis_name"`         // Зв'язок з виконавцем
			Number          string    `json:"number"`           // Номер виконавчого провадження
			Category        string    `json:"category"`         // Категорія стягнення
			Id              string    `json:"id"`               // Ідентіфікаційний номер
			DepartmentPhone string    `json:"department_phone"` // Номер телефону виконавця
			Executor        string    `json:"executor"`         // Виконавець
			ExecutorPhone   string    `json:"executor_phone"`   // Номер телефону виконавця
			ExecutorEmail   string    `json:"executor_email"`   // Email виконавця
			DeductionType   string    `json:"deduction_type"`   // Категорія стягнення
			LastName        string    `json:"last_name"`        // Прізвище боржника
			FirstName       string    `json:"first_name"`       // Ім'я боржника
			MiddleName      string    `json:"middle_name"`      // Ім'я по батькові боржника
			BirthDate       time.Time `json:"birth_date"`       // Дата народження боржника
		} `json:"items"`
	} `json:"data"`
}

// GetPenalties
// Отримання інформації про актуальні виконавчі провадження приватної особи за ПІБ
// https://docs.opendatabot.com/#/%D0%92%D0%B8%D0%BA%D0%BE%D0%BD%D0%B0%D0%B2%D1%87%D1%96%20%D0%BF%D1%80%D0%BE%D0%B2%D0%B0%D0%B4%D0%B6%D0%B5%D0%BD%D0%BD%D1%8F/penaltiesByFioAndBirth
func (odb *OdbClient) GetPenalties(
	firstName string, // Ім’я боржника
	lastName string, // Прізвище боржника
	birthDate string, // Дата народження у форматі YYYY-MM-DD
	params map[string]string, // map[string]string{
	//	"middle_name":		"По-батькові боржника",
	//	"categories[1]":	"Код категорії", //01 - стягнення коштів
	//  //02 - звернення стягнення на майно
	//  //03 - стягнення аліментів
	//  //04 - стягнення періодичних платежів (крім аліментів)
	//  //05 - стягнення заборгованості із заробітної плати та інших платежів, пов’язаних з трудовими відносинами
	//  //06 - стягнення соціальних виплат
	//  //07 - стягнення заборгованості з оплати комунальних послуг
	//  //08 - стягнення штрафів у справах про адміністративні правопорушення
	//  //09 - стягнення штрафів у справах про адміністративні правопорушення у сфері безпеки дорожнього руху
	//  //10 - забезпечення позову
	//  //11 - зобов’язання вчинити певні дії або утриматися від їх вчинення
	//  //12 - поновлення на роботі
	//  //13 - вселення стягувача
	//  //14 - виселення
	//  //15 - відібрання дитини
	//  //16 - заборона вчиняти певні дії
	//  //17 - конфіскація майна
	//  //18 - конфіскація майна, вилученого митними органами
	//  //18.1 - конфіскація майна засуджених
	//  //19 - конфіскація коштів та майна за вчинення корупційного та пов’язаного з корупцією правопорушення
	//  //20 - оплатне вилучення
	//  //21 - передача стягувачу предметів, зазначених у виконавчому документі
	//  //22 - стягнення коштів на користь держави
	//  //23 - рішення Європейського суду з прав людини
	//  //24 - стягнення виконавчого збору
	//  //25 - стягнення витрат виконавчого провадження
	//  //26 - стягнення штрафів, накладених державним, приватним виконавцем
	//  //27 - стягнення основної винагороди приватного виконавця
	//  //28 - усунення перешкод у побаченні з дитиною, встановлення побачення з дитиною
	//}
) (response *PenaltyByFioSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	params["first_name"] = firstName
	params["last_name"] = lastName
	params["birth_date"] = birthDate

	err = odb.Do(penaltiesEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 0,
	//    "items": [
	//      {
	//        "court_name": "суд загальної юрисдикції Селидівський міський суд Донецької області.Суддя:І.М.Владимирська",
	//        "gis_name": "Селидівський міський відділ державної виконавчої служби Головного територіального управління юстиції у Донецькій області",
	//        "number": "57654570",
	//        "category": "стягнення коштів на користь держави",
	//        "id": "12570284040",
	//        "department_phone": "(06237) 9-08-09",
	//        "executor": "Глухівський міськрайонний відділ державної виконавчої служби",
	//        "executor_phone": "(06237) 9-08-09",
	//        "executor_email": "example@gmail.com",
	//        "deduction_type": "стягнення аліментів",
	//        "last_name": "Ткаченко",
	//        "first_name": "Олег",
	//        "middle_name": "Васильович",
	//        "birth_date": "1975-05-06T00:00:00Z"
	//      }
	//    ]
	//  }
	//}
}

type KoatuuRegions struct {
	Status string `json:"status"` // Статус запиту
	Data   []struct {
		Code string `json:"code"` // Код об'єкту
		Name string `json:"name"` // Назва об'єкту
		Type string `json:"type"` // Тип об'єкту
	} `json:"data"`
}

// GetKoatuuRegions
// Отримати список всіх областей України
// https://docs.opendatabot.com/#/%D0%9A%D0%9E%D0%90%D0%A2%D0%A3%D0%A3/koatuuRegions
func (odb *OdbClient) GetKoatuuRegions() (response *KoatuuRegions, err error) {
	err = odb.Do(koatuuRegionsEndpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": [
	//    {
	//      "code": "0100000000",
	//      "name": "Автономна Республіка Крим",
	//      "type": "region"
	//    },
	//    {
	//      "code": "0500000000",
	//      "name": "Вінницька",
	//      "type": "region"
	//    },
	//    {
	//      "code": "0700000000",
	//      "name": "Волинська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "1200000000",
	//      "name": "Дніпропетровська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "1400000000",
	//      "name": "Донецька",
	//      "type": "region"
	//    },
	//    {
	//      "code": "1800000000",
	//      "name": "Житомирська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "2100000000",
	//      "name": "Закарпатська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "2300000000",
	//      "name": "Запорізька",
	//      "type": "region"
	//    },
	//    {
	//      "code": "2600000000",
	//      "name": "Івано-Франківська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "3200000000",
	//      "name": "Київська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "3500000000",
	//      "name": "Кіровоградська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "4400000000",
	//      "name": "Луганська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "4600000000",
	//      "name": "Львівська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "4800000000",
	//      "name": "Миколаївська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "5100000000",
	//      "name": "Одеська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "5300000000",
	//      "name": "Полтавська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "5600000000",
	//      "name": "Рівненська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "5900000000",
	//      "name": "Сумська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "6100000000",
	//      "name": "Тернопільська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "6300000000",
	//      "name": "Харківська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "6500000000",
	//      "name": "Херсонська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "6800000000",
	//      "name": "Хмельницька",
	//      "type": "region"
	//    },
	//    {
	//      "code": "7100000000",
	//      "name": "Черкаська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "7300000000",
	//      "name": "Чернівецька",
	//      "type": "region"
	//    },
	//    {
	//      "code": "7400000000",
	//      "name": "Чернігівська",
	//      "type": "region"
	//    },
	//    {
	//      "code": "8000000000",
	//      "name": "М.київ",
	//      "type": "region"
	//    },
	//    {
	//      "code": "8500000000",
	//      "name": "М.севастополь",
	//      "type": "region"
	//    }
	//  ]
	//}
}

type Koatuu struct {
	Status string `json:"status"`
	Data   struct {
		Name  string `json:"name"`
		Code  string `json:"code"`
		Type  string `json:"type"`
		Items struct {
			RegionDistrict []struct {
				Code string `json:"code"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"region-district"`
			CityAndDistrict []struct {
				Code string `json:"code"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"city-and-district"`
			City []struct {
				Code      string `json:"code"`
				Name      string `json:"name"`
				Type      string `json:"type"`
				Districts []struct {
					Code string `json:"code"`
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"districts"`
			} `json:"city"`
		} `json:"items"`
	} `json:"data"`
}

// GetKoatuuRegionsByCode
// Отримати список всіх округів, районів та міст за КОАТУУ кодом
// https://docs.opendatabot.com/#/%D0%9A%D0%9E%D0%90%D0%A2%D0%A3%D0%A3/koatuu
func (odb *OdbClient) GetKoatuuRegionsByCode(
	code string, // КОАТУУ код (10 або 17 цифр)
) (response *Koatuu, err error) {
	endpoint := fmt.Sprintf(koatuuRegionsByCodeEndpoint, code)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "name": "Одеська",
	//    "code": "5100000000",
	//    "type": "region",
	//    "items": {
	//      "city": [
	//        {
	//          "name": "Одеса",
	//          "code": "5110100000",
	//          "type": "city",
	//          "districts": [
	//            {
	//              "name": "Київський",
	//              "code": "5110136900",
	//              "type": "city-district"
	//            },
	//            {
	//              "name": "Малиновський",
	//              "code": "5110137300",
	//              "type": "city-district"
	//            },
	//            {
	//              "name": "Приморський",
	//              "code": "5110137500",
	//              "type": "city-district"
	//            },
	//            {
	//              "name": "Суворовський",
	//              "code": "5110137600",
	//              "type": "city-district"
	//            }
	//          ]
	//        }
	//      ],
	//      "city-and-district": [
	//        {
	//          "name": "Балта",
	//          "code": "5110200000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Білгород-Дністровський",
	//          "code": "5110300000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Біляївка",
	//          "code": "5110500000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Ізмаїл",
	//          "code": "5110600000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Чорноморськ",
	//          "code": "5110800000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Подільськ",
	//          "code": "5111200000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Теплодар",
	//          "code": "5111500000",
	//          "type": "city-and-district"
	//        },
	//        {
	//          "name": "Южне",
	//          "code": "5111700000",
	//          "type": "city-and-district"
	//        }
	//      ],
	//      "region-district": [
	//        {
	//          "name": "Ананьївський",
	//          "code": "5120200000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Арцизький",
	//          "code": "5120400000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Балтський",
	//          "code": "5120600000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Білгород-Дністровський",
	//          "code": "5120800000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Біляївський",
	//          "code": "5121000000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Березівський",
	//          "code": "5121200000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Болградський",
	//          "code": "5121400000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Великомихайлівський",
	//          "code": "5121600000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Іванівський",
	//          "code": "5121800000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Ізмаїльський",
	//          "code": "5122000000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Кілійський",
	//          "code": "5122300000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Кодимський",
	//          "code": "5122500000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Лиманський",
	//          "code": "5122700000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Подільський",
	//          "code": "5122900000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Окнянський",
	//          "code": "5123100000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Любашівський",
	//          "code": "5123300000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Миколаївський",
	//          "code": "5123500000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Овідіопольський",
	//          "code": "5123700000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Роздільнянський",
	//          "code": "5123900000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Ренійський",
	//          "code": "5124100000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Савранський",
	//          "code": "5124300000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Саратський",
	//          "code": "5124500000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Тарутинський",
	//          "code": "5124700000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Татарбунарський",
	//          "code": "5125000000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Захарівський",
	//          "code": "5125200000",
	//          "type": "region-district"
	//        },
	//        {
	//          "name": "Ширяївський",
	//          "code": "5125400000",
	//          "type": "region-district"
	//        }
	//      ]
	//    }
	//  }
	//}
}

type RealtySuccess struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		Count          string `json:"count"` // Кількість знайдениї об'єктів нерухомості
		ReportResultId string `json:"reportResultId"`
		Items          []struct {
			DcGroupType string `json:"dcGroupType"`
			Name        string `json:"name"` // Адреса нерухомості
			Id          string `json:"id"`   // ID об'єкту нерухомості
			Link        string `json:"link"` // Посилання на повний об'єкт нерухомості
		} `json:"items"`
	} `json:"data"`
}

// GetRealty
// Отримати інформацію щодо всіх об’єктів нерухомості, земельних ділянок або обтяжень по компанії або фізичній особі
// https://docs.opendatabot.com/#/%D0%9D%D0%B5%D1%80%D1%83%D1%85%D0%BE%D0%BC%D1%96%D1%81%D1%82%D1%8C/realty
func (odb *OdbClient) GetRealty(
	code string, // код ЄДРПОУ або ІПН
	params map[string]string, //map[string]string{
	//	"offset":	"Зміщення відносно початку результатів пошуку",
	//	"limit":	"Кількість записів",
	//	"timeout":	"Кількість секунд очікування відповіді від реєстру майнових прав",
	//	"role":		"Роль суб’єкта", //3 - Обтяжувач
	//	//4 - Особа, майно/права якої обтяжуються
	//	//6 - Іпотекодержатель
	//	//7 - Майновий поручитель
	//	//8 - Іпотекодавець
	//	//9 - Боржник
	//	//10 - Особа, в інтересах якої встановлено обтяження
	//	//11 - Власник
	//	//12 - Правонабувач
	//	//13 - Правокористувач
	//	//14 - Землевласник
	//	//15 - Землеволоділець
	//	//16 - Інший
	//	//17 - Наймач
	//	//18 - Орендар
	//	//19 - Наймодавець
	//	//20 - Орендодавець
	//	//21 - Управитель
	//	//22 - Вигодонабувач
	//	//23 - Установник
	//	//25 - Довірчій власник
	//}
) (response *RealtySuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	params["code"] = code

	err = odb.Do(realtyEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": "1",
	//    "reportResultId": "5001395171432",
	//    "items": [
	//      {
	//        "dcGroupType": "1",
	//        "name": "Київська обл., м. Київ, вулиця Дзержинського, будинок 71",
	//        "id": "45035236",
	//        "link": "https://opendatabot.com/api/v2/realty/5001396491688/45268910?apiKey=xxxxxxxxxx"
	//      }
	//    ]
	//  }
	//}
}

type RealtyItemSuccess struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		ResultId         string `json:"resultId"`           // ID результату
		ObjectResultLink string `json:"object_result_link"` // Посилання на отримання результату про отримання повної інформації про нерухомість
	} `json:"data"`
}

// GetRealtyById
// Замовити витяг з докладною інформацією по об’єкту нерухомості або земельній ділянці
// https://docs.opendatabot.com/#/%D0%9D%D0%B5%D1%80%D1%83%D1%85%D0%BE%D0%BC%D1%96%D1%81%D1%82%D1%8C/realty-object
func (odb *OdbClient) GetRealtyById(
	reportResultId string, // Ідентифікатор групи адресів суб'єкта
	id string, // Ідентифікатор об'єкта групи reportResultId
) (response *RealtyItemSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(realtyByIdEndpoint, reportResultId, id)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "resultId": "057557bde3148f33a3d787c615e9404b",
	//    "object_result_link": "https://opendatabot.com/api/v2/realty-result?activity_id=057557bde3148f33a3d787c615e9404b&apiKey=xxxxxxxxxx"
	//  }
	//}
}

type RealtyResultSuccess struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		Data struct {
			Realty            string `json:"realty"`            // Актуальна інформація про нерухоміть
			OldMortgageJson   string `json:"oldMortgageJson"`   // Інформація про іпотеку(до 2013р)
			OldLimitationJson string `json:"oldLimitationJson"` // Інформація про обтяження(до 2013р)
			OldRealty         string `json:"oldRealty"`         // Інформація про нерухомість(до 2013р)
			AllAdresses       string `json:"allAdresses"`       // Інші адреса
		} `json:"data"`
		Status  string `json:"status"`   // Статус обробки запиту
		PdfLink string `json:"pdf_link"` // Посилання на PDF документ
		Fixed   string `json:"fixed"`    // Статус виправлення PDF документу
	} `json:"data"`
}

// GetRealtyResult
// Отримання витягу або поточного статусу його формування по об’єкту нерухомості або земельній ділянці
// https://docs.opendatabot.com/#/%D0%9D%D0%B5%D1%80%D1%83%D1%85%D0%BE%D0%BC%D1%96%D1%81%D1%82%D1%8C/realty-result
func (odb *OdbClient) GetRealtyResult(
	resultId string, // Ідентифікатор пошуку за результатом витягу
) (response *RealtyResultSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(realtyResultEndpoint, map[string]string{
		"resultId": resultId,
	}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "data": {
	//      "realty": "",
	//      "oldMortgageJson": "",
	//      "oldLimitationJson": "",
	//      "oldRealty": "",
	//      "allAdresses": ""
	//    },
	//    "status": "pending",
	//    "pdf_link": "https://opendatabot.com/pdf/realty/106/4a092ddc72544c2ab233dd4aa0aca77a-106745--cc354e103ff85c7c8a24dfd0b25f1312.pdf",
	//    "fixed": "false"
	//  }
	//}
}

type RealtyObjectReportSuccess struct {
	Status string `json:"status"` // Статус запиту
	Data   struct {
		ResultId         string `json:"resultId"`           // ID результату
		ObjectResultLink string `json:"object_result_link"` // Посилання на отримання результату про отримання повної інформації про нерухомість
	} `json:"data"`
}

// GetRealtyReportByNumber
// Замовити витяг з докладною інформацією по об’єкту нерухомості або земельній ділянці за кадастровим номером або за номером реєстрації
// https://docs.opendatabot.com/#/%D0%9D%D0%B5%D1%80%D1%83%D1%85%D0%BE%D0%BC%D1%96%D1%81%D1%82%D1%8C/realty-object-report
func (odb *OdbClient) GetRealtyReportByNumber(
	number string, // кадастровий номер (XXXXXXXXXX:XX:XXX:XXXX) або код реєстрації (максімально 28 цифр)
) (response *RealtyObjectReportSuccess, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(realtyReportByNumberEndpoint, number)

	err = odb.Do(endpoint, map[string]string{}, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "resultId": "871982275b1080d392fe404f8b04fe8b",
	//    "object_result_link": "https://opendatabot.ua/api/v2/realty-result?activity_id=871982275b1080d392fe404f8b04fe8b&apiKey=xxxxxxxxxx"
	//  }
	//}
}

type Timeline struct {
	Status string `json:"status"`
	Data   struct {
		Count int `json:"count"`
		Items []struct {
			LogId     string    `json:"log_id"`
			Id        string    `json:"id"`
			Code      string    `json:"code"`
			Type      string    `json:"type"`
			CreatedAt time.Time `json:"created_at"`
			EventDate time.Time `json:"event_date"`
			Change    []struct {
				OldValue          string   `json:"old_value,omitempty"`
				NewValue          string   `json:"new_value,omitempty"`
				Number            string   `json:"number,omitempty"`
				DocumentId        string   `json:"document_id,omitempty"`
				CountAddedItems   string   `json:"countAddedItems,omitempty"`
				AddedItems        []string `json:"addedItems,omitempty"`
				CountRemovedItems string   `json:"countRemovedItems,omitempty"`
				RemovedItems      string   `json:"removedItems,omitempty"`
				Date              string   `json:"date,omitempty"`
				Name              string   `json:"name,omitempty"`
				IsCompany         string   `json:"is_company,omitempty"`
				JudgmentCode      string   `json:"judgment_code,omitempty"`
				Source            string   `json:"source,omitempty"`
				Link              string   `json:"link,omitempty"`
				CompanyName       string   `json:"company_name,omitempty"`
				WithoutChangeLogs string   `json:"without_change_logs,omitempty"`
				DeclarantId       string   `json:"declarant_id,omitempty"`
				Year              string   `json:"year,omitempty"`
				DeclarationId     string   `json:"declaration_id,omitempty"`
				PublicType        string   `json:"public_type,omitempty"`
				SubjectType       string   `json:"subject_type,omitempty"`
				CodePdv           string   `json:"code_pdv,omitempty"`
				EventDate         string   `json:"eventDate,omitempty"`
				StartDate         string   `json:"startDate,omitempty"`
				EndDate           string   `json:"endDate,omitempty"`
				Termless          string   `json:"termless,omitempty"`
				SanctionList      string   `json:"sanctionList,omitempty"`
				SanctionReason    string   `json:"sanctionReason,omitempty"`
				Pib               string   `json:"pib,omitempty"`
				Resident          string   `json:"resident,omitempty"`
			} `json:"change"`
		} `json:"items"`
	} `json:"data"`
}

// GetTimeline
// Отримання стрічки змін за реєстрами
// https://docs.opendatabot.com/#/%D0%9C%D0%BE%D0%BD%D1%96%D1%82%D0%BE%D1%80%D0%B8%D0%BD%D0%B3%20%D0%B1%D1%96%D0%B7%D0%BD%D0%B5%D1%81%D1%83/timeline
func (odb *OdbClient) GetTimeline(
	params map[string]string, // map[string]string{
	//	"code":			"код ЄДРПОУ",
	//	"from_id":		"Зміщення відносно log_id",
	//	"type":			"Кількість записів", // change_status_borrower - зміна статусу виконавчого провадження у якості боржника
	//	//change_status_creditor - зміна статусу виконавчого провадження у якості стягувача
	//	//new_penalty_borrower - нове виконавче провадження у якості боржника
	//	//new_penalty_creditor - нове виконавче провадження у якості стягувача
	//	//penalty - нове виконавче провадження в реєстрі боржників
	//	//realty - зміна об'єктів нерухомості у реєстрі речових прав
	//	//wagedebt - нова заборгованість по виплаті заробітної плати
	//	//inspections - нова перевірка контролюючими органами
	//	//debt - зміна статусу податкового боргу
	//	//new_court_defendant - новий судовий процес у якості відповідача
	//	//add_court_defendant - додано нового відповідача по вже існуючій справі
	//	//new_court_plaintiff - новий судовий процес у якості позивача
	//	//add_court_plaintiff - додано нового позивача по вже існуючій справі
	//	//new_court_third_person - новий судовий процес у якості третьої сторони
	//	//add_court_third_person - додано третю сторону по вже існуючій справі
	//	//new_decision - новий документ за судовою справою
	//	//new_schedule - нове засідання у судовій справі
	//	//legal - реєстраційні зміни компанії
	//	//legal_declarant - власник компанії є декларантом
	//	//edr_company - реєстраційні зміни компанії (архівні події)
	//	//bankruptcy_fop - Інформація щодо банкрутства ФОП
	//	//bankruptcy_company - Інформація щодо банкрутства юридичних осіб
	//	//bankruptcy_person - Інформація щодо банкрутства фізичних осіб
	//	//beneficiaries_user - зміни власників компанії
	//	//vat - наявність у компанії свідоцтва платника ПДВ
	//	//drorm - Інформація по обтяженням рухомого майна
	//	//sanction - Санкція юридичної особи
	//	//person_sanction - Санкція фізичної особи
	//	"pib":			"Прізвище, ім'я, по батькові (тільки для типу person_sanction)",
	//	"itn":			"ІНН (тільки для типу person_sanction)",
	//	"date_start":	"Фільтр за датою початку події (event_date) у форматі Y-m-d",
	//	"date_end":		"Фільтр за датою закінчення події (event_date) у форматі Y-m-d",
	//	"created_date":	"Фільтр за датою створення (created_date) у форматі Y-m-d",
	//	"offset":		"Зміщення відносно початку результатів пошуку",
	//	"limit":		"Кількість записів",
	//	"order":		"Порядок сортування. Available values : asc, desc",
	//	"order_field":	"Поле сортування. Available values : id, created_at, event_date",
	//}
) (response *Timeline, err error) {
	if err = checkApiKey(odb); err != nil {
		return nil, err
	}

	err = odb.Do(timelineEndpoint, params, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
	//{
	//  "status": "ok",
	//  "data": {
	//    "count": 1,
	//    "items": [
	//      {
	//        "log_id": "123",
	//        "id": "2CLH936",
	//        "code": "12345678",
	//        "type": "penalty",
	//        "created_at": "2018-05-18T00:00:00+03:00",
	//        "event_date": "2019-04-24T00:00:00+03:00",
	//        "change": [
	//          {
	//            "old_value": "Відкрито",
	//            "new_value": "Примусове виконання",
	//            "number": "59004247",
	//            "document_id": "34217403"
	//          },
	//          {
	//            "old_value": "Відкрито",
	//            "new_value": "Примусове виконання",
	//            "number": "59004247",
	//            "document_id": "34217403"
	//          },
	//          {
	//            "number": "59004247",
	//            "document_id": "34217403"
	//          },
	//          {
	//            "number": "59004247",
	//            "document_id": "34217403"
	//          },
	//          {
	//            "document_id": "03342184",
	//            "countAddedItems": "1",
	//            "addedItems": [
	//              "59003451"
	//            ],
	//            "countRemovedItems": "2",
	//            "removedItems": "[59003454, 59003455]"
	//          },
	//          {
	//            "document_id": "2926320991",
	//            "countAddedItems": "1",
	//            "addedItems": [
	//              "Харківська обл., Харківський р., м. Мерефа, вулиця Жуковського, будинок 54, квартира 454"
	//            ],
	//            "countRemovedItems": "0",
	//            "removedItems": "[ ]"
	//          },
	//          {
	//            "old_value": "18181.02",
	//            "new_value": "443.67",
	//            "date": "2022-04-14"
	//          },
	//          {
	//            "old_value": "Заплановано",
	//            "new_value": "Проведено",
	//            "document_id": "2250460"
	//          },
	//          {
	//            "old_value": "160",
	//            "new_value": "16073",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "is_company": "true"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "b17fdb43b2164dfe5be068ab5462cd2eb17fdb43b2164dfe5be068ab5462cd2e",
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "document_id": "94928370",
	//            "source": "decision",
	//            "number": "904/9488/21",
	//            "judgment_code": "3",
	//            "link": "https://opendatabot.com/api/v2/court/94928370?apiKey="
	//          },
	//          {
	//            "document_id": "c72cd1248badd590e86f74e94540e410",
	//            "source": "schedule",
	//            "number": "904/9488/21",
	//            "judgment_code": "3"
	//          },
	//          {
	//            "old_value": "зареєстровано",
	//            "new_value": "порушено справу про банкрутство",
	//            "company_name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ЮПІТЕР",
	//            "without_change_logs": "false"
	//          },
	//          {
	//            "declarant_id": "1182296",
	//            "old_value": "Куліш Олександр Сергійович",
	//            "new_value": "Куліш Олександр Сергійович",
	//            "year": "2017",
	//            "declaration_id": "101a4e5e-c4db-4cfb-9363-247646781ea2"
	//          },
	//          {
	//            "company_name": "АО АДВОКАТ ГЛОБАЛ",
	//            "old_value": "Україна, 65045",
	//            "new_value": "Україна, 65045, Одеська обл., місто Одеса, вул.Новосельського"
	//          },
	//          {
	//            "name": "Фізична особа - підприємець Машев Віталій Валерійович",
	//            "public_type": "Повідомлення про поновлення провадження у справі про банкрутство з визнанням мирової угоди недійсною або її розірвання"
	//          },
	//          {
	//            "name": "Товариство з обмеженою відповідальністю Династія",
	//            "public_type": "Оголошення про порушення справи про банкрутство"
	//          },
	//          {
	//            "name": "Слав Галина Дмитрівна",
	//            "public_type": "Оголошення про порушення справи про банкрутство"
	//          },
	//          {
	//            "company_name": "ОК ЗАТИШНИЙ",
	//            "old_value": "СКАДОВСЬКА МІСЬКА РАДА СКАДОВСЬКОГО РАЙОНУ ХЕРСОНСЬКОЇ ОБЛАСТІ",
	//            "new_value": "СКАДОВСЬКА МІСЬКА РАДА"
	//          },
	//          {
	//            "name": "ТОВАРИСТВО З ОБМЕЖЕНОЮ ВІДПОВІДАЛЬНІСТЮ ГІДРОПНЕВМОАГРЕГАТ",
	//            "subject_type": "company",
	//            "code_pdv": "224963515089",
	//            "eventDate": "2022-02-18 23:22:11",
	//            "old_value": "anul",
	//            "new_value": "payer"
	//          },
	//          {
	//            "document_id": "32563846",
	//            "countAddedItems": "1",
	//            "addedItems": [
	//              "16798865"
	//            ],
	//            "countRemovedItems": "0",
	//            "removedItems": "[ ]"
	//          },
	//          {
	//            "name": "Товариство з обмеженою відповідальністю Натекс",
	//            "startDate": "2021-05-22",
	//            "endDate": "2024-05-22",
	//            "termless": "false",
	//            "sanctionList": "РНБО",
	//            "sanctionReason": "1086/2004"
	//          },
	//          {
	//            "pib": "Mohamad Nouman Ala Al Din",
	//            "startDate": "2021-05-22",
	//            "endDate": "2024-05-22",
	//            "termless": "false",
	//            "sanctionList": "РНБО",
	//            "sanctionReason": "1086/2004",
	//            "resident": "false"
	//          }
	//        ]
	//      }
	//    ]
	//  }
	//}
}

func checkApiKey(odb *OdbClient) error {
	if odb.Settings.ApiKey == "" {
		return errors.New("ApiKey is not specified")
	}

	return nil
}

func checkNotEmpty(id string) error {
	if id == "" {
		return errors.New("Id is not specified")
	}

	return nil
}

func buildQueryParams(endpoint string, params map[string]string) (uri string, err error) {
	base, err := url.Parse(endpoint)

	if err != nil {
		return
	}

	query := url.Values{}

	for key, value := range params {
		query.Add(key, value)
	}

	base.RawQuery = query.Encode()

	return base.String(), err
}

// Do
// Make Request
func (odb *OdbClient) Do(endpoint string, params map[string]string, v interface{}) (err error) {
	if odb.Settings.ApiKey != "" {
		params["apiKey"] = odb.Settings.ApiKey
	}

	endpointWithParams, err := buildQueryParams(endpoint, params)

	if err != nil {
		return err
	}

	resp, err := http.Get(endpointWithParams)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &v)

	if err != nil {
		return err
	}

	return nil
}
