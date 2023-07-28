package model

type AllocationStatus string

const (
	ALLOCATED   AllocationStatus = "allocated"
	CHECKED_IN  AllocationStatus = "checked_in"
	CHECKED_OUT AllocationStatus = "checked_out"
	STAYED      AllocationStatus = "stayed"
	NO_SHOW     AllocationStatus = "no_show"
)

//class Status extends Enum {
//
//    const ALLOCATED = 'allocated';
//    const CHECKED_IN = 'checked_in';
//    const CHECKED_OUT = 'checked_out';
//    const STAYED = 'stayed';//@deprecated
//    const NO_SHOW = 'no_show';//MAXIM: internal use with future perspective do it changeable in WEB UI
//
//    protected static $sSequence = [
//        self::ALLOCATED => self::CHECKED_IN,
//        self::CHECKED_IN => self::CHECKED_OUT,
//        self::CHECKED_OUT =>  self::STAYED,//self::STAYED,
//        self::STAYED => null//@deprecated
//    ];
//
//    public function nextStatusInSequence() {
//        $next = self::$sSequence[$this->_value];
//        return is_null($next) ? null : new self($next);
//    }
//
//    protected static function initializeValues() {
//        return [
//            self::ALLOCATED => __('Allocated'),
//            self::CHECKED_IN => __('Check-In'),
//            self::STAYED => __('Stayed'),//@deprecated
//            self::CHECKED_OUT => __('Check-Out'),
//            self::NO_SHOW => __('No Show'),
//        ];
//    }
//
//    public function asDataArray($options = []) {
//        $includeNext = true;
//        if (array_key_exists('next', $options)) {
//            $includeNext = $options['next'];
//        }
//
//        $data = parent::asDataArray();
//
//        if ($includeNext) {
//            $next = $this->nextStatusInSequence();
//            if (is_null($next)) {
//                $data['next'] = $next;
//            } else
//                {
//                $data['next'] = $next->asDataArray(['next' => false]);
//            }
//        }
//
//        return $data;
//    }
//
//}
