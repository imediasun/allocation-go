package model

//type Reservation struct {
//	ID                int
//	Creator           Agent
//	Price             ValueObjectMoney
//	CreationDate      time.Time
//	Status            ReservationStatus
//	ProviderReference string
//	Channel           string
//	Remark            string
//	Client            string
//	Manual            bool
//	PaymentOption     ValueObjectReservationPaymentOption
//	Groups            []ReservationGroup
//	CancellationDate  time.Time
//	CheckoutDate      time.Time
//	CheckoutTime      time.Time
//	ArrivalDate       time.Time
//	ArrivalTime       time.Time
//	DepartureDate     time.Time
//	DepartureTime     time.Time
//	Adults            int
//	Children          int
//	Infants           int
//	Source            string
//	SourceReference   string
//	SourceReference2  string
//	SourceReference3  string
//}

//class Reservation extends UnitOfWork\Entity\LifeCycleEventListener implements \JsonSerializable, PaidEntityInterface
//{
//
//    use jSerializeTrait;
//
//    /** @var null|int */
//    private $id = null;
//
//    /** @var Agent */
//    private $creator = null;
//
//    /** @var ValueObject\Money */
//    private $price = null;
//
//    /** @var \DateTime */
//    private $creationDate;
//
//    /** @var ValueObject\Reservation\Status */
//    private $status = null;
//
//    /**
//     * Reservation ID from Channel service
//     * @var null|string
//     */
//    private $providerReference = null;
//
//    /** @var null|string */
//    private $channel = null;
//
//    /** @var null|string */
//    private $remark = null;
//
//    /** @var null|string */
//    private $client = null;
//
//    /** @var bool */
//    private $manual = false;
//
//    /** @var null|ValueObject\Reservation\PaymentOption */
//    private $paymentOption = null;
//
//    /** @var ReservationGroup[] */
//    private $groups;
//
//    /** @var null|\DateTime */
//    private $cancellationDate = null;
//
//    /** @var null|\DateTime */
//    private $_startDate = null;
//
//    /** @var null|\DateTime */
//    private $_endDate = null;
//
//    /** @var SegmentReservation */
//    private $segment = null;
//
//    /** @var ValueObject\Service|null */
//    private $source = null;
//
//    /** @var null|mixed $logs */
//    private $logs = null;
//
//    /** @var CurrencyRate[]|null */
//    private $currencyRates = null;
//
//    /** @var bool */
//    private $foct = false;
//
//    /** @var bool */
//    private $isCityTaxToProvider = false;
//
//    /** @var null|int */
//    private $metaGroupID = null;
//
//    /** @var null|Client */
//    private $customer = null;
//
//    /** @var Color|null */
//    private $color = null;
//
//    /** @var bool|null */
//    private $isVirtualCC = null;
//
//    /** @var string|null */
//    private $checkInLink = null;
//
//    /** @var Company[]|EntityCollection */
//    private $companies;
//
//    /** @var ValueObject\Reservation\CommissionType|null */
//    private $channelCommissionType = null;
//
//    /** @var float|null */
//    private $channelCommissionValue = null;
//
//    /** @var float|null */
//    private $releazeTime = null;
//
//    /** @var ValueObject\Time|null */
//    private $requestedCheckInTime = null;
//
//
//    /** @var ValueObject\Time|null */
//    private $requestedCheckOutTime = null;
//
//    const DEFAULT_SEGMENT_ID = 1;
//
//    /**
//     * Reservation constructor.
//     * @param ValueObject\Money $price
//     * @param ValueObject\Reservation\Status|null $status
//     * @param \DateTime|null $creationDate
//     * @param false $isManual
//     * @param null $providerReference
//     * @param null $channel
//     * @param null $remark
//     * @param null $client
//     * @param ValueObject\Reservation\PaymentOption|null $paymentOption
//     * @param ValueObject\Service|null $reservationSource
//     * @param null $billingRates
//     * @param Color|null $color
//     * @param bool $isVirtualCC
//     */
//    public function __construct(
//        ValueObject\Money $price,
//        ValueObject\Reservation\Status $status = null,
//        \DateTime $creationDate = null,
//        $isManual = false,
//        $providerReference = null,
//        $channel = null,
//        $remark = null,
//        $client = null,
//        ValueObject\Reservation\PaymentOption $paymentOption = null,
//        ValueObject\Service $reservationSource = null,
//        $billingRates = null,
//        Color $color = null,
//        bool $isVirtualCC = false
//    )
//    {
//        parent::__constructor();
//
//        $this->groups = [];
//
//        $this->creationDate = is_null($creationDate) ? Moment::getNow() : $creationDate;
//        $this->status = $status;
//        $this->changePrice($price);
//        $this->manual = $isManual;
//        $this->providerReference = $providerReference;//MAXIM: here is not correct name. This is subChannel of $this->source
//        $this->channel = $channel;
//        $this->remark = $remark;
//        $this->client = $client;
//        $this->paymentOption = $paymentOption;
//        $this->source = $reservationSource;//MAXIM: here is not correct name. This is Channel::XXX or must be 'new Channel' object
//        $this->billingRates = $billingRates;
//        $this->color = $color;
//        $this->isVirtualCC = $isVirtualCC;
//    }
//
//    /**
//     * @return bool
//     */
//    public function canBeModified(): bool
//    {
//        try {
//            /** @TODO  */
//        } catch (\Exception $exception) {
//            return true;
//        }
//
//        return true;
//    }
//
//    /**
//     * @return void
//     * @throws Exception
//     */
//    public function throwIfCannotBeModified()
//    {
//        if(!$this->canBeModified()) {
//            throw new \Exception("Reservation can't be modified!");
//        }
//    }
//
//    /**
//     * @return \string[][]
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_attachChannelAsCompany()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_persistDependedEntities()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_validate()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_onBeforePersist()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_validateProducts()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_checkAllocations()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_inventoryCache()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_updateBillingDocs()
//     * @uses \Zenith\Application\Domain\Entity\Reservation::_onFetched()
//     */
//    protected function _getSubscribedEvents()
//    {
//        return [
//            UnitOfWork\Entity\LifeCycleEventListener::BEFORE_COMMIT => ['_attachChannelAsCompany','_persistDependedEntities'],
//            UnitOfWork\Entity\LifeCycleEventListener::BEFORE_CREATE => ['_validate', '_onBeforePersist', '_validateProducts', '_markVCC'],
//            UnitOfWork\Entity\LifeCycleEventListener::BEFORE_UPDATE => ['_validate', '_validateProducts', '_checkAllocations', '_markVCC'],
//            UnitOfWork\Entity\LifeCycleEventListener::AFTER_UPDATE => ['_inventoryCache', '_updateBillingDocs'],
//            UnitOfWork\Entity\LifeCycleEventListener::AFTER_CREATE => ['_inventoryCache', '_updateBillingDocs'],
//            UnitOfWork\Entity\LifeCycleEventListener::ON_FETCHED => ['_onFetched']
//        ];
//    }
//
//    /**
//     * @return int|null
//     */
//    public function id(): ?int
//    {
//        return $this->id;
//    }
//
//    /**
//     * @return \DateTime
//     * @throws Exception
//     */
//    public function startDate()
//    {
//        $s = strtotime('+10 years');
//
//        foreach ($this->groups() as $g) {
//            $sd = $g->startDate()->getTimestamp();
//
//            if ($sd < $s) $s = $sd;
//        }
//
//        $this->_startDate = Moment::getNow();
//        $this->_startDate->setTimestamp($s);
//
//        if (is_null($this->_startDate)) throw new \LogicException('Start date was not initialized');
//        return $this->_startDate;
//    }
//
//    /**
//     * @return \DateTime
//     * @throws Exception
//     */
//    public function endDate()
//    {
//        $e = strtotime('-10 years');
//
//        foreach ($this->groups() as $g) {
//            $ed = $g->endDate()->getTimestamp();
//
//            if ($ed > $e) $e = $ed;
//        }
//
//        $this->_endDate = Moment::getNow();
//        $this->_endDate->setTimestamp($e);
//
//        if (is_null($this->_endDate)) throw new \LogicException('End date was not initialized');
//        return $this->_endDate;
//    }
//
//    /**
//     * @return Agent|null
//     */
//    public function creator()
//    {
//        return $this->creator;
//    }
//
//    /**
//     * @return ValueObject\Money
//     */
//    public function price()
//    {
//        return $this->price;
//    }
//
//    /**
//     * @return \DateTime
//     */
//    public function createdAt()
//    {
//        return $this->creationDate;
//    }
//
//    /**
//     * @return ValueObject\Reservation\Status
//     */
//    public function status()
//    {
//        return $this->status;
//    }
//
//    /**
//     * @return null|string
//     */
//    public function providerReference()
//    {
//        return $this->providerReference;
//    }
//
//    /** @return null|ValueObject\Service */
//    public function source()
//    {
//        return $this->source;
//    }
//
//    /** @return null|mixed */
//    public function logs()
//    {
//        return $this->logs;
//    }
//
//    /**
//     * @return boolean
//     */
//    public function isFoct()
//    {
//        return $this->foct;
//    }
//
//    /** @return null|int */
//    public function metaGroupID(): ?int
//    {
//        return $this->metaGroupID;
//    }
//
//    /**
//     * @return Client|null
//     */
//    public function customer(): ?Client
//    {
//        return $this->customer;
//    }
//
//    /**
//     * @return int
//     */
//    public function paxCount() : int
//    {
//        $result = 0;
//        foreach ($this->groups() as $group) {
//            foreach ($group->pax() as $client) {
//                $result++;
//            }
//        }
//        return $result;
//    }
//
//    /**
//     * @param Client $client
//     * @return Reservation
//     */
//    public function changeCustomer(Client $client): self
//    {
//        $this->customer = $client;
//        return $this;
//    }
//
//    /**
//     * @return string|null
//     */
//    public function color(): ?string
//    {
//        return $this->color;
//    }
//
//    /**
//     * @return null|string
//     */
//    public function getDisplayProviderReference()
//    {
//        if (is_null($this->providerReference)) return null;
//        if (strpos($this->providerReference, '-')) {
//            $data = explode('-', $this->providerReference, 2);
//        } else {
//            $data[1] = $this->providerReference;
//        }
//        return @$data[1] ?: null;
//    }
//
//    /**
//     * @return null|string
//     */
//    public function channel()
//    {
//        return $this->channel;
//    }
//
//    /**
//     * @return null|string
//     */
//    public function remark()
//    {
//        return $this->remark;
//    }
//
//    /**
//     * @return boolean
//     */
//    public function isManual()
//    {
//        return $this->manual;
//    }
//
//    /**
//     * @return \DateTime|null
//     */
//    public function cancellationDate()
//    {
//        return $this->cancellationDate;
//    }
//
//    /**
//     * @return null|ValueObject\Reservation\PaymentOption
//     */
//    public function paymentOption()
//    {
//        return $this->paymentOption;
//    }
//
//    /**
//     * @return SegmentReservation|null;
//     */
//    public function segment()
//    {
//        return $this->segment;
//    }
//
//    /**
//     * @return CurrencyRate[]|null
//     */
//    public function currencyRates()
//    {
//        return $this->currencyRates;
//    }
//
//    /**
//     * @param int $infants
//     * @param int $children
//     * @param int $adults
//     * @param int $seniors
//     * @param int $nights
//     */
//    public function paxStatisticCollect(&$infants, &$children, &$adults, &$seniors, &$nights)
//    {
//        foreach ($this->groups() as $group) {
//            $nights += $group->durationInDays() - 1;
//            $group->paxStatisticCollect($infants, $children, $adults, $seniors);
//        }
//    }
//
//    /**
//     * Assign creator to reservation
//     * @param Agent $creator
//     * @param bool $force
//     * @throws ViolationException
//     */
//    public function assignCreator(Agent $creator, $force = false)
//    {
//        if (!is_null($this->creator) && !$force) throw new ViolationException('Reservation already assigned to a creator');
//        $this->creator = $creator;
//    }
//
//    /**
//     * @param Company $company
//     */
//    public function assignCompany(Company $company) {
//        if(is_null($this->companies)) {
//            $this->companies = new EntityCollection();
//        }
//
//        $this->companies->add((string)$company->id(), $company);
//    }
//
//    /**
//     * @return Company[]|EntityCollection
//     */
//    public function companies()
//    {
//        return !is_null($this->companies) ? $this->companies : new EntityCollection();
//    }
//
//    /**
//     * @param ValueObject\Reservation\PaymentOption $option
//     */
//    public function _changePaymentOption(ValueObject\Reservation\PaymentOption $option)
//    {
//        $this->paymentOption = $option;
//    }
//
//    /**
//     * @return null|string
//     */
//    public function client()
//    {
//        return $this->client;
//    }
//
//    /**
//     * @param \DateTime $startDate
//     * @param \DateTime $endDate
//     * @param Country $paxNationality
//     * @return ReservationGroup Added group
//     */
//    public function addGroup(\DateTime $startDate, \DateTime $endDate, Country $paxNationality = null)
//    {
//        $group = new ReservationGroup($startDate, $endDate, $paxNationality);
//        $group->_notifyReservationCurrencyChanged($this->price->currency());
//        $this->groups[] = $group;
//        return $group;
//    }
//
//    /**
//     * @param int $id
//     * @return bool True on success, false on failure
//     */
//    public function removeGroupById($id)
//    {
//        foreach ($this->groups as $gIndex => $g) {
//            if ($g->id() === $id) {
//                unset($this->groups[$gIndex]);
//                return true;
//            }
//        }
//
//        return false;
//    }
//
//    /**
//     * @param int $id
//     * @return null|ReservationGroup
//     */
//    public function getGroupById($id)
//    {
//        foreach ($this->groups as $g) {
//            if ($g->id() === $id) {
//                return $g;
//            }
//        }
//        return null;
//    }
//
//    /**
//     * @param $parentId
//     * @return null|ReservationGroup
//     */
//    public function getGroupByParentId($parentId): ?ReservationGroup
//    {
//        foreach ($this->groups as $g) {
//            if ($g->parentId() === $parentId) {
//                return $g;
//            }
//        }
//        return null;
//    }
//
//    public function getItemByProductId(string $productId): ?Reservation\Item
//    {
//        foreach ($this->groups() as $g) {
//            foreach ($g->items() as $i) {
//                if ($i->product()->id()->equals(new \Zenith\Application\Domain\Identifier\Product($productId))) {
//                    return $i;
//                }
//            }
//        }
//        return null;
//    }
//
//    public function getGroupByItemId($id)
//    {
//        foreach ($this->groups() as $g) {
//            if (!is_null($g->getItemById($id))) {
//                return $g;
//            }
//        }
//        return null;
//    }
//
//    /**
//     * @param $id
//     * @return Reservation\Item|null
//     */
//    public function getItemById($id)
//    {
//        foreach ($this->groups() as $g) {
//            $item = $g->getItemById($id);
//            if (!is_null($item))
//                return $item;
//        }
//        return null;
//    }
//
//    /**
//     * @return string
//     */
//    public function leadName(): string
//    {
//        $lead = $this->lead();
//        return !is_null($lead) ?
//            $lead->fullName() :
//            '';
//    }
//
//    /**
//     * @return int|null
//     */
//    public function leadId(): ?int
//    {
//        $lead = $this->lead();
//        return !is_null($lead) ?
//            $lead->id() :
//            null;
//    }
//
//    /**
//     * @return Client|null
//     */
//    public function lead(): ?Client
//    {
//        $anyGroup = $this->getFirstGroup();
//        return !is_null($anyGroup) ?
//            (current($anyGroup->pax()) ?: null) :
//            null;
//    }
//
//    /**
//     * @return ReservationGroup|null
//     */
//    public function getFirstGroup()
//    {
//        foreach ($this->groups() as $g) {
//            return $g;
//        }
//        return null;
//    }
//
//    /**
//     * Return first found Venue ID
//     * @return int|null
//     */
//    public function venueId()
//    {
//        $venueId = null;
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if (
//                    $item->type()->equals(ValueObject\Reservation\Item\Type::PRODUCT) &&
//                    !is_null($item->product()) &&
//                    !is_null($item->product()->venueId())
//                ) {
//                    $venueId = $item->product()->venueId();
//                    break 2;
//                }
//            }
//        }
//
//        return $venueId;
//    }
//
//    /**
//     * @return ReservationGroup[]
//     */
//    public function groups()
//    {
//        return $this->groups;
//    }
//
//    public function getRoomNumbers()
//    {
//        $roomNumbers = [];
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if (
//                    !is_null($item->product()) &&
//                    !is_null($item->allocatedObject()) &&
//                    !is_null($item->allocatedObject()->productObject())
//                ) {
//                    $roomNumbers[] = $item->allocatedObject()->productObject()->name();
//                }
//            }
//        }
//        return $roomNumbers;
//    }
//
//    /**
//     * @param ValueObject\Money $price
//     */
//    public function changePrice(ValueObject\Money $price)
//    {
//        $this->price = $price;
//        foreach ($this->groups as $group) {
//            foreach ($group->items() as $item) {
//                $item->_notifyBookingCurrencyChanged($this->price()->currency());
//            }
//        }
//    }
//
//    /**
//     * @param ValueObject\Reservation\Status $status
//     * @return bool
//     */
//    public function isStatus(ValueObject\Reservation\Status $status)
//    {
//        return $this->status->equals($status);
//    }
//
//    public function amendAllocationStatus(ValueObject\Allocation\Status $status)
//    {
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if (is_null($item->allocatedObject()))
//                    $item->changeAllocationObject(
//                        new ValueObject\Reservation\Item\AllocatedObject($status)
//                    );
//                else
//                    $item->allocatedObject()->changeStatus($status);
//            }
//        }
//    }
//
//    /**
//     * @param ValueObject\Reservation\Status $status
//     * @throws Exception
//     */
//    public function changeStatus(ValueObject\Reservation\Status $status)
//    {
////        if ($this->status->equals($status)) {
////            return; //do nothing
////        }
////
//        // this is needed for cancellationDate of booking, don't comment out this whoever you are
//        if ($status->equals(ValueObject\Reservation\Status::CANCELLED)) {
//            $this->cancellationDate = Moment::getNow();
////            $this->_updateAvailability(false);
//        }
////
////        if ($status->equals(ValueObject\Reservation\Status::CONFIRMED)){
////            $this->_updateAvailability(true);
////        }
//
//        $this->status = $status;
//    }
//
////    /**
////     * @param \DateTime $date
////     */
////    public function changeCreationDate(\DateTime $date) {
////        if ($date->getTimestamp() === $this->creationDate->getTimestamp()) {
////            return; //do nothing
////        }
////
////        $this->creationDate = $date;
////    }
//
//    /**
//     * @param string|null $reference
//     */
//    public function changeProviderReference($reference)
//    {
//        $this->providerReference = $reference;
//    }
//
//    /**
//     * @param string|null $channel
//     */
//    public function changeChannel($channel)
//    {
//        $this->channel = $channel;
//    }
//
//    /**
//     * @param string|null $remark
//     */
//    public function changeRemark($remark = null)
//    {
//        $this->remark = $remark;
//    }
//
//    /**
//     * @param SegmentReservation|null $segment
//     */
//    public function changeSegment(SegmentReservation $segment = null)
//    {
//        $this->segment = $segment;
//    }
//
//    /**
//     * @param boolean $foct
//     */
//    public function setFoct($foct)
//    {
//        $this->foct = $foct;
//    }
//
//    /**
//     * @param Agent $agent
//     * @return bool
//     * @throws IncompleteEntityException
//     */
//    public function canBeViewedBy(Agent $agent)
//    {
//        if (is_null($this->creator)) throw new IncompleteEntityException('Reservation does not have a creator assigned to it');
//
//        if ($this->creator->account()->id() == $agent->account()->id()) return true;
//        if ($this->creator->isSameAs($agent)) return true;
//
//        foreach ($this->groups as $group) {
//            if ($group->canBeViewedBy($agent)) return true;
//        }
//
//        return false;
//    }
//
//    /**
//     * @param false $ignoreChannel
//     * @param bool $ignorePermissions
//     * @param bool $onlyProviderName
//     * @return string    The PmsProvider must not return because default args of $this->getProvider
//     */
//    public function getDisplayProvider($ignoreChannel = false, $ignorePermissions = true, $onlyProviderName = true)
//    {
//        return $this->getProvider($ignoreChannel, $ignorePermissions, $onlyProviderName);
//    }
//
//    /**
//     * @param bool $ignoreChannel
//     * @param bool $ignorePermissions
//     * @param bool $onlyProviderName
//     * @param bool $ignoreSource
//     * @return string|PmsProvider
//     */
//    function getProvider($ignoreChannel = false, $ignorePermissions = false, $onlyProviderName = true, $ignoreSource = true)
//    {
//        if ((is_null($this->channel) && $ignoreSource) || $ignoreChannel) {
//            return implode(' / ', [$this->creator->account()->name(), $this->creator->name()]);
//        } else {
//            //@TODO horrible! HORRIBLE
//            $venueID = $this->venueId();
//            if (is_null($venueID))
//                return '';
//
//            /** @var Venue $venue */
//            $venue = EntityManager::create(true)
//                ->getQuery(Venue::class)
//                ->with(Venues::PMS_DATA)
//                ->findOneById($venueID);
//            if (is_null($venue))
//                return '';
//
//            /** @var PmsProvider|null $pmsData */
//            $pmsData = null;
//            $sourceData = null;
//            $sourceId = is_null($this->source) ? null : $this->source->id();
//            /** @var PmsProvider $pd */
//            foreach ($venue->pmsProviderData() ?: [] as $pd) {
//                if ($pd->provider()->code() === $this->channel()) {
//                    $pmsData = $pd;
//                    break;
//                }
//                if ($sourceId === $pd->provider()->name()) {
//                    $sourceData = $pd;
//                }
//            }
//
//            if (is_null($pmsData) && !is_null($sourceData)) $pmsData = $sourceData;
//
//            if (is_null($pmsData))
//                return '';
//
//            if ($ignorePermissions) {
//                return $onlyProviderName ? $pmsData->provider()->name() : $pmsData;
//            }
//
//            return $pmsData->isManagedByHotel() ? $pmsData->provider()->name() : 'Hotel Connect';
//        }
//    }
//
//    /**
//     * @return int[]
//     */
//    public function getInvolvedAccountIds()
//    {
//        $accountIds = [];
//
//        foreach ($this->groups as $g) {
//            foreach ($g->items() as $i) {
//                $id = $i->accountId();
//                if (!is_null($accountIds)) {
//                    $accountIds[] = $id;
//                }
//            }
//        }
//
////        $accountIds[] = $this->creator->account()->id();
//        return array_unique($accountIds);
//    }
//
//    /**
//     * @return ValueObject\Money
//     * @throws Exception
//     */
//    public function getBalance() : ValueObject\Money
//    {
//        $balance = Accounting::create()->getBalance(new Refers(Refers::RESERVATION), $this->id);
//        return new ValueObject\Money($balance['total']['leftToPay'], $balance['total']['currency']);
//    }
//
//    /**
//     * @return bool
//     */
//    public function canBeAmended()
//    {
//        //New PMS
//        return $this->creationDate->getTimestamp() > strtotime('2016-03-12');
//    }
//
//    protected function _attachChannelAsCompany(UnitOfWork $uow, ?ChangeSet $changes, ?self $self){
//        //check if reservation source is the Service object and company is not assigned to reservation
//        if($this->source() instanceof ValueObject\Service && $this->companies()->count() === 0) {
//            //Check if the service is a channel
//            if(is_subclass_of($this->source()->config()->get('API'), AbstractChannelApi::class)) {
//                //Fetch a company that assigned with channel
//                $companyCriteria = new \Zenith\Application\Infrastructure\ObjectQuery\Criteria\Company();
//                $companyCriteria->changeChannel($this->source());
//                /** @var Company $company */
//                $company = $uow->getQuery(Company::class)->fetchOne($companyCriteria);
//                //If we found company that assigned with channel
//                // then attach the company to the reservation
//                // and attach the customer to the company
//                if(!is_null($company)) {
//                    $this->assignCompany($company);
//                    if($this->customer() instanceof Client) {
//                        $company->addCustomers([$this->customer()]);
//                    }
//                }
//
//            }
//        }
//    }
//
//    //region Event listeners
//    protected function _persistDependedEntities(UnitOfWork $uow)
//    {
//        $register = function ($entity) use ($uow) {
//            if (!$uow->isRegistered($entity))
//                $uow->register($entity);
//        };
//        $register($this->customer);
//        foreach ($this->groups() as $g) {
//            foreach ($g->pax() as $p)
//                $register($p);
//        }
//    }
//
//    protected function _onBeforePersist()
//    {
//        if (is_null($this->status)) {
//            $total = 0;
//            $confirmed = 0;
//            $failed = 0;
//            $onRequest = 0;
//
//            $confirmedStatus = new ValueObject\Reservation\Item\Status(ValueObject\Reservation\Item\Status::CONFIRMED);
//            $failedStatus = new ValueObject\Reservation\Item\Status(ValueObject\Reservation\Item\Status::FAILED);
//            $onRequestStatus = new ValueObject\Reservation\Item\Status(ValueObject\Reservation\Item\Status::ON_REQUEST);
//
//            foreach ($this->groups() as $g) {
//                foreach ($g->items() as $i) {
//                    $total++;
//                    if ($i->status()->equals($confirmedStatus)) $confirmed++;
//                    elseif ($i->status()->equals($onRequestStatus)) $onRequest++;
//                    elseif ($i->status()->equals($failedStatus)) $failed++;
//                }
//            }
//
//            if ($onRequest > 0) $this->status = new ValueObject\Reservation\Status(ValueObject\Reservation\Status::ON_REQUEST);
//            else if ($confirmed === $total) $this->status = new ValueObject\Reservation\Status(ValueObject\Reservation\Status::CONFIRMED);
//            elseif ($failed === $total) $this->status = new ValueObject\Reservation\Status(ValueObject\Reservation\Status::FAILED);
//            else $this->status = new ValueObject\Reservation\Status(ValueObject\Reservation\Status::PARTIALLY_CONFIRMED);
//        }
//
//    }
//
//    protected function _checkAllocations()
//    {
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if (
//                    in_array($item->status()->value(), [ValueObject\Reservation\Item\Status::CANCELLED, ValueObject\Reservation\Item\Status::FAILED]) ||
//                    $this->status()->equals(ValueObject\Reservation\Status::CANCELLED)
//                )
//                    $item->changeAllocationObject(null);
//            }
//        }
//    }
//
//    /**
//     * @return bool
//     */
//    public function isImmediatelyNonRefundable()
//    {
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if ($item->isImmediatelyNonRefundable()) {
//                    return true;
//                }
//            }
//        }
//        return false;
//    }
//
//    public function isNonRefundable()
//    {
//        foreach ($this->groups() as $group) {
//            foreach ($group->items() as $item) {
//                if ($item->isNonRefundable($this->creationDate->diff($group->startDate())->days)) {
//                    return true;
//                }
//            }
//        }
//        return false;
//    }
//
//    /**
//     * @throws InvalidEntityException
//     */
//    protected function _validate()
//    {
//        if (is_null($this->creator)) throw new InvalidEntityException('Reservation does not have a creator');
////        if (count($this->groups) === 0) throw new InvalidEntityException('No groups found on reservation');//MAXIM: Disable it test because channels (like BookingCom) returning canceled reservations without Groups, But we need to store it reservations
//
//        foreach ($this->groups() as $g) {
//            $items = $g->items();
//            $pax = $g->pax();
//
//            if (count($items) === 0) throw new InvalidEntityException('Group with empty items found on reservation');
//            if (count($pax) === 0) throw new InvalidEntityException('Group with empty pax found on reservation');
//        }
//    }
//
//    /**
//     * @throws InvalidEntityException
//     */
//    protected function _validateProducts()
//    {
//        foreach ($this->groups() as $g) {
//            foreach ($g->items() as $ik => $iv) {
//                $iv->validate($ik);
//            }
//        }
//    }
//
//    /**
//     * @throws InvalidEntityException
//     */
//    protected function _markVCC()
//    {
//        foreach($this->customer()->creditCards() as $creditCard){
//            $creditCard->markIfVirtualCC();
//            if($creditCard->isVirtualCC()){
//                $this->isVirtualCC = true;
//            }
//        }
//    }
//
//    /**
//     * @throws Exception
//     */
//    protected function _onFetched()
//    {
//        foreach ($this->groups() as $g) {
//            $g->_notifyOnFetched();
//        }
//    }
//
//    //endregion
//
//    /**
//     * @param UnitOfWork $uow
//     * @param ChangeSet|null $changeSet
//     * @param Reservation|null $oldThis
//     */
//    protected function _inventoryCache(UnitOfWork $uow, ChangeSet $changeSet = null, Reservation $oldThis = null)
//    {
//        $inventoryService = Inventory::create();
//        /** @var self $r */
//        foreach ([$oldThis, $this] as $r) {
//            if (is_null($r)) continue;
//
//            foreach ($r->groups() as $g) {
//                $start = $g->startDate();
//                $end = $g->endDate();
//                foreach ($g->items() as $i) {
//                    if (
//                        !is_null($i->product()) &&
//                        $i->product()->type()->equals(ValueObject\Product\Type::ROOM)
//                    ) {
//                        $inventoryService->regenerate(
//                            $start,
//                            $end,
//                            $i->product()->id()
//                        );
//                    }
//                }
//            }
//        }
//    }
////MAXIM: I did not find use (useless)
////    /**
////     * @return array
////     */
////    public function getCalculatedAllocationStatuses()
////    {
////        $today = Moment::getNow();
////        $calculatedStatuses = [];
////        foreach ($this->groups() as $group) {
////            foreach ($group->items() as $item) {
////                if (is_null($item->product()) || !$item->product()->type()->equals(ValueObject\Product\Type::ROOM)) continue;
////                $allocationObject = $item->allocatedObject();
////                $allocatedStatus = is_null($allocationObject) ? null : $allocationObject->status();
////                $calculatedStatus = null;
////                if (
////                    is_null($item->product()) ||
////                    !$item->product()
////                        ->type()
////                        ->equals(\Zenith\Application\Domain\ValueObject\Product\Type::ROOM)
////                ) continue;
////
////                if ($group->endDate() <= $today) {
////                    $calculatedStatus = ValueObject\Allocation\AllocationCalculatedStatus::NO_SHOW;
////                    if (
////                        !is_null($allocatedStatus) &&
////                        (
////                            $allocatedStatus->equals(ValueObject\Allocation\Status::CHECKED_OUT) ||
////                            $allocatedStatus->equals(ValueObject\Allocation\Status::STAYED) ||
////                            $allocatedStatus->equals(ValueObject\Allocation\Status::CHECKED_IN)
////                        )
////                    ) {
////                        $calculatedStatus =
////                            $group->endDate() == $today
////                            && $allocatedStatus->equals(ValueObject\Allocation\Status::CHECKED_IN)
////                                ? ValueObject\Allocation\AllocationCalculatedStatus::DEPARTURE
////                                : ValueObject\Allocation\AllocationCalculatedStatus::STAYED;
////                    }
////                } elseif (
////                    $group->startDate() == $today &&
////                    (
////                        is_null($allocatedStatus) ||
////                        $allocatedStatus->equals(ValueObject\Allocation\Status::ALLOCATED)
////                    )
////                ) {
////                    $calculatedStatus = ValueObject\Allocation\AllocationCalculatedStatus::ARRIVAL;
////                } elseif ($group->startDate() > $today) {
////                    $calculatedStatus = ValueObject\Allocation\AllocationCalculatedStatus::FUTURE;
////                } else {//in house
////                    $calculatedStatus =
////                        !is_null($allocatedStatus)
////                        && $allocatedStatus->equals(ValueObject\Allocation\Status::CHECKED_IN)
////                            ? ValueObject\Allocation\AllocationCalculatedStatus::IN_HOUSE
////                            : ValueObject\Allocation\AllocationCalculatedStatus::NO_SHOW;
////                }
////
////                $calculatedStatuses[] = [
////                    'room_id' => $item->product()->id()->__toString(),
////                    'allocated_object_id' => $item->allocatedObject() ? $item->allocatedObject()->productObject()->id() : '',
////                    'state' => $calculatedStatus,
////                    'group_id' => $group->id(),
////                ];
////            }
////        }
////
////        return $calculatedStatuses;
////    }
//
//    /**
//     * @param UnitOfWork|null $uow
//     * @param ChangeSet|null $changeSet
//     * @throws Exception
//     */
//    public function _updateBillingDocs(UnitOfWork $uow = null, ChangeSet $changeSet = null)
//    {
//        try {
//            //do not update or create documents if the group end date is in the past
//            if ($this->endDate()->setTime(23, 59, 59) < Moment::getNow()->setTime(23, 59, 59)) return;
//
//            $accounting = Accounting::create();
//
//            $oldFoctValue = null;
//            if (!is_null($changeSet) && $changeSet->hasChangeFor('foct')) {
//                $oldFoctValue = $changeSet->getChangeFor('foct')->getOriginalValue();
//            }
//
//            foreach ($this->groups as $group) {
//                //defining default payers
//                $payers = $accounting->figureOutReservationPayers($this, reset($this->groups));
//                $accounting->cancelIrrelevantDocuments($payers, $this);
//
//                //if documents with an appropriate group id parameter and status of "active" exist, then we need to update them. Otherwise to create new ones
//                if ($accounting->groupDocumentsExist($group->id())) {
//                    $addItem = false;
//                    $addedItems = [];
//
//                    //this loop is meant to add new document items as a standalone. It happens when a user adds new products using "Add Item" form
//                    foreach ($group->items() as $item) {
//                        if (!is_null($item->getPayer())) {
//                            if (!$addItem) $addItem = true;
//
//                            $payer = $item->getPayer()->equals(ValueObject\Accounting\FromToEntity::CUSTOMER)
//                                || $item->getPayer()->equals(ValueObject\Accounting\FromToEntity::GROUP)
//                                ? new ValueObject\Accounting\FromToEntity(ValueObject\Accounting\FromToEntity::CUSTOMER)
//                                : $item->getPayer();
//
//                                if ($item->getPayer()->equals(ValueObject\Accounting\FromToEntity::CUSTOMER)
//                                    || $item->getPayer()->equals(ValueObject\Accounting\FromToEntity::CORPORATE)) {
//                                    $payerId = $item->getPayerId();
//                                } elseif ($item->getPayer()->equals(ValueObject\Accounting\FromToEntity::GROUP)) {
//                                    $gp = $group->pax();
//                                    $lead = reset($gp);
//                                    $payerId = $lead->id();
//                                } else $payerId = null;
//                                //here is the check if the item should be added as a standalone
//                                if (!in_array($item->id(), $addedItems) && !is_null($changeSet)) {
//                                    $accounting->checkIfStandalone($item, $this->id, $payer->value(), $group->id(), $payerId, $changeSet);
//                                };
//                                $addedItems[] = $item->id();
//                            }
//                        }
//                        //if an item was not added as a standalone, then the system will update documents and put new items on its own
//                        if (!$addItem) {
//                            $accounting->updateSystemBillingDocs($this, $group->id(), $oldFoctValue, $changeSet);
//                        }
//                }
//                else {
//                    //create new group documents
//                    $accounting->createSystemBillingDocs($this, $group);
//                }
//            }
//        } catch (Exception $e) {
//            if(is_null($uow)) {
//                throw $e;
//            }
//        }
//    }
//
//    /**
//     * @param Color $color
//     * @return $this
//     */
//    public function changeColor(Color $color): Reservation
//    {
//        $this->color = $color;
//
//        return $this;
//    }
//
//    /**
//     * @param string|null $url
//     * @return $this
//     */
//    public function setCheckInLink(?string $url = null): self
//    {
//        $this->checkInLink = $url;
//        return $this;
//    }
//
//    /**
//     * @return string|null
//     */
//    public function checkInLink(): ?string
//    {
//        return $this->checkInLink;
//    }
//
//
//    /**
//     * @param boolean $status
//     * @return Reservation
//     */
//    public function setVirtualCC(bool $status): Reservation
//    {
//        $this->isVirtualCC = $status;
//        return $this;
//    }
//
//    /**
//     * @return bool
//     */
//    public function isVirtualCC(): bool
//    {
//        return $this->isVirtualCC;
//    }
//
//    /**
//     * @param ValueObject\Reservation\CommissionType $ct
//     */
//    public function setCommissionType(ValueObject\Reservation\CommissionType $ct)
//    {
//        $this->channelCommissionType = $ct;
//    }
//
//
//    /**
//     * @param $cv
//     */
//    public function setCommissionValue($cv)
//    {
//        $this->channelCommissionValue = $cv;
//    }
//
//    /**
//     * @param $rt
//     */
//    public function setReleazeTime($rt)
//    {
//        $this->releazeTime = $rt;
//    }
//
//    /**
//     * @return int|null
//     */
//    public function releazeTime(): ?int {
//        return $this->releazeTime;
//    }
//
//    /**
//     * Return Booking Cancellation Status according to date, if date is not provided current date-time will be used
//     *
//     * @param \DateTime|null $date
//     * @return string
//     */
//    public function getBookingCancellationPolicyStatus(?\DateTime $date = null) : string {
//        $statuses = [];
//        /**
//         * Loop through the groups and get cancellation status each of them
//         */
//        foreach ($this->groups() as $group) {
//            array_push($statuses, $group->getGroupCancellationPolicyStatus($date));
//        }
//        $statuses = array_unique($statuses);
//        /**
//         * if we have only one unique status then this is the status of the booking
//         * if we have 2 or 3 difference statuses we have 100% probability that will be at least one option where client can money back, but not all, so partially-refundable
//         * otherwise non-refundable (we can't determinate)
//         */
//        switch (count($statuses)) {
//            case 1: return array_pop($statuses);
//            case 2:
//            case 3: return 'partial';
//            default: return 'nonrefundable';
//        }
//    }
//
//    /**
//     * @param $rt
//     */
//    public function setRequestedCheckInTime($rt)
//    {
//        if($rt instanceof ValueObject\Time){
//            $this->requestedCheckInTime = $rt;
//        }else if(!empty($rt)) {
//            $this->requestedCheckInTime = ValueObject\Time::createFromHis($rt);
//        }
//    }
//
//    /**
//     * @return ValueObject\Time|null
//     */
//    public function requestedCheckInTime(): ?ValueObject\Time {
//        return $this->requestedCheckInTime;
//    }
//
//
//    /**
//     * @param $rt
//     */
//    public function setRequestedCheckOutTime($rt)
//    {
//        if($rt instanceof ValueObject\Time){
//            $this->requestedCheckOutTime = $rt;
//        }else if(!empty($rt)) {
//            $this->requestedCheckOutTime = ValueObject\Time::createFromHis($rt);
//        }
//    }
//
//    /**
//     * @return ValueObject\Time|null
//     */
//    public function requestedCheckOutTime(): ?ValueObject\Time {
//        return $this->requestedCheckOutTime;
//    }
//
//    /**
//     * @param int|string $id
//     * @return Company|null
//     */
//    public function getCompanyById($id): ?Company
//    {
//        return $this->companies()->get($id);
//    }
//
//    /**
//     * @param \SelectQuery $selectQuery
//     * @return void
//     */
//    public function modifyDocumentCriteria(\SelectQuery $selectQuery)
//    {
//        $selectQuery->where('ad.ReservationID', $this->id());
//    }
//
//    /**
//     * @return int|null
//     */
//    public function getCustomerId(): ?int
//    {
//        return $this->customer() ? $this->customer()->id() : null;
//    }
//}
