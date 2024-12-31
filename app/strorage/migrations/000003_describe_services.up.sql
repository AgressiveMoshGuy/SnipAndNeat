CREATE TABLE services_names (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT
);



INSERT INTO services_names (id, name, description)
VALUES
    (1, 'MarketplaceNotDeliveredCostItem', 'Возврат невостребованного товара от покупателя на склад.'),
    (2, 'MarketplaceReturnAfterDeliveryCostItem', 'Возврат от покупателя на склад после доставки.'),
    (3, 'MarketplaceDeliveryCostItem', 'Доставка товара до покупателя.'),
    (4, 'MarketplaceSaleReviewsItem', 'Приобретение отзывов на платформе.'),
    (5, 'ItemAdvertisementForSupplierLogistic', 'Доставка товаров на склад Ozon — кросс-докинг.'),
    (6, 'OperationMarketplaceServiceStorage', 'Размещения товаров.'),
    (7, 'MarketplaceMarketingActionCostItem', 'Продвижение товаров.'),
    (8, 'MarketplaceServiceItemInstallment', 'Продвижениe и продажа в рассрочку.'),
    (9, 'MarketplaceServiceItemMarkingItems', 'Обязательная маркировка товаров.'),
    (10, 'MarketplaceServiceItemFlexiblePaymentSchedule', 'Гибкий график выплат.'),
    (11, 'MarketplaceServiceItemReturnFromStock', 'Комплектация товаров для вывоза продавцом.'),
    (12, 'ItemAdvertisementForSupplierLogisticSeller', 'Транспортно-экспедиционные услуги.'),
    (13, 'ItemAgentServiceStarsMembership', 'Вознаграждение за услугу «Звёздные товары».'),
    (14, 'MarketplaceServiceItemDelivToCustomer', 'Последняя миля.'),
    (15, 'MarketplaceServiceItemDirectFlowTrans', 'Магистраль.'),
    (16, 'MarketplaceServiceItemDropoffFF', 'Обработка отправления.'),
    (17, 'MarketplaceServiceItemDropoffPVZ', 'Обработка отправления.'),
    (18, 'MarketplaceServiceItemDropoffSC', 'Обработка отправления.'),
    (19, 'MarketplaceServiceItemFulfillment', 'Сборка заказа.'),
    (20, 'MarketplaceServiceItemPickup', 'Выезд транспортного средства по адресу продавца для забора отправлений — Pick-up.'),
    (21, 'MarketplaceServiceItemReturnAfterDelivToCustomer', 'Обработка возврата.'),
    (22, 'MarketplaceServiceItemReturnFlowTrans', 'Обратная магистраль.'),
    (23, 'MarketplaceServiceItemReturnNotDelivToCustomer', 'Обработка отмен.'),
    (24, 'MarketplaceServiceItemReturnPartGoodsCustomer', 'Обработка невыкупа.'),
    (25, 'MarketplaceRedistributionOfAcquiringOperation', 'Оплата эквайринга.'),
    (26, 'MarketplaceReturnStorageServiceAtThePickupPointFbsItem', 'Краткосрочное размещение возврата FBS.'),
    (27, 'MarketplaceReturnStorageServiceInTheWarehouseFbsItem', 'Долгосрочное размещение возврата FBS.'),
    (28, 'MarketplaceServiceItemDeliveryKGT', 'Доставка крупногабаритного товара (КГТ).'),
    (29, 'MarketplaceServiceItemDirectFlowLogistic', 'Логистика.'),
    (30, 'MarketplaceServiceItemReturnFlowLogistic', 'Обратная логистика.'),
    (31, 'MarketplaceServicePremiumCashbackIndividualPoints', 'Услуга продвижения «Бонусы продавца».'),
    (32, 'MarketplaceServicePremiumPromotion', 'Услуга продвижение Premium, фиксированная комиссия.'),
    (33, 'OperationMarketplaceWithHoldingForUndeliverableGoods', 'Удержание за недовложение товара.'),
    (34, 'MarketplaceServiceItemDropoffPPZ', 'Услуга drop-off в пункте приёма заказов.'),
    (35, 'MarketplaceServiceItemRedistributionReturnsPVZ', 'Перевыставление возвратов на ПВЗ.'),
    (36, 'OperationMarketplaceAgencyFeeAggregator3PLGlobal', 'Тарификация агентской услуги Agregator 3PL Global.'),
    (37, 'MarketplaceServiceItemDirectFlowLogisticVDC', 'Логистика вРЦ.');