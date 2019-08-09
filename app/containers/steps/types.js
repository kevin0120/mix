import PropTypes from 'prop-types';

export type tStepProps = {
  bindAction: (PropTypes.Element)=>mixed,
  step: {},
  isCurrent: boolean,
  bindDescription: (PropTypes.Element)=>mixed
};
